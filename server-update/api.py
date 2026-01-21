#!/usr/bin/env python3
"""
Servidor de Atualizações para 00cli

Este servidor serve a API de atualizações para o 00cli.
Quando o binário é compilado com 'make build', ele automaticamente
atualiza o servidor com a nova versão.

Uso:
    python api.py                    # Inicia na porta 8080
    python api.py --port 9000        # Inicia na porta 9000
    python api.py --host 0.0.0.0     # Aceita conexões externas

Endpoints:
    GET /latest         - Retorna info da última versão
    GET /download/<bin> - Download do binário específico
    GET /health         - Health check
    POST /upload        - Upload de novo binário (via make build)
"""

import os
import sys
import json
import hashlib
import argparse
import subprocess
from pathlib import Path
from datetime import datetime
from typing import Optional, Dict, Any

try:
    from flask import Flask, jsonify, send_file, request, abort
except ImportError:
    print("❌ Flask não instalado. Execute:")
    print("   pip install flask")
    print("   ou")
    print("   pip install -r requirements.txt")
    sys.exit(1)

# Configurações
BASE_DIR = Path(__file__).parent.absolute()
BINARIES_DIR = BASE_DIR / "binaries"
METADATA_FILE = BASE_DIR / "metadata.json"
PROJECT_ROOT = BASE_DIR.parent

# Garantir que o diretório de binários existe
BINARIES_DIR.mkdir(exist_ok=True)

app = Flask(__name__)

# Plataformas suportadas
PLATFORMS = {
    "linux-amd64": {"goos": "linux", "goarch": "amd64", "ext": ""},
    "linux-arm64": {"goos": "linux", "goarch": "arm64", "ext": ""},
    "darwin-amd64": {"goos": "darwin", "goarch": "amd64", "ext": ""},
    "darwin-arm64": {"goos": "darwin", "goarch": "arm64", "ext": ""},
    "windows-amd64": {"goos": "windows", "goarch": "amd64", "ext": ".exe"},
}


def get_current_version() -> str:
    """Obtém a versão atual do projeto"""
    try:
        # Tentar via git tag
        result = subprocess.run(
            ["git", "describe", "--tags", "--always", "--dirty"],
            cwd=PROJECT_ROOT,
            capture_output=True,
            text=True
        )
        if result.returncode == 0 and result.stdout.strip():
            return result.stdout.strip()
    except Exception:
        pass

    return "v0.1.0"


def get_file_hash(filepath: Path) -> str:
    """Calcula SHA256 do arquivo"""
    sha256 = hashlib.sha256()
    with open(filepath, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            sha256.update(chunk)
    return sha256.hexdigest()


def load_metadata() -> Dict[str, Any]:
    """Carrega metadata das releases"""
    if METADATA_FILE.exists():
        with open(METADATA_FILE, "r") as f:
            return json.load(f)
    return {
        "tag_name": get_current_version(),
        "name": f"Release {get_current_version()}",
        "published_at": datetime.utcnow().isoformat() + "Z",
        "body": "Atualização automática",
        "assets": []
    }


def save_metadata(metadata: Dict[str, Any]):
    """Salva metadata das releases"""
    with open(METADATA_FILE, "w") as f:
        json.dump(metadata, f, indent=2)


def scan_binaries(base_url: str) -> list:
    """Escaneia binários disponíveis"""
    assets = []
    for filename in BINARIES_DIR.iterdir():
        if filename.is_file() and filename.name.startswith("00cli"):
            stat = filename.stat()
            assets.append({
                "name": filename.name,
                "browser_download_url": f"{base_url}/download/{filename.name}",
                "size": stat.st_size,
                "sha256": get_file_hash(filename)
            })
    return assets


def build_binary(platform: str) -> Optional[Path]:
    """Compila binário para uma plataforma específica"""
    if platform not in PLATFORMS:
        return None

    config = PLATFORMS[platform]
    version = get_current_version()

    binary_name = f"00cli-{platform}{config['ext']}"
    output_path = BINARIES_DIR / binary_name

    env = os.environ.copy()
    env["GOOS"] = config["goos"]
    env["GOARCH"] = config["goarch"]
    env["CGO_ENABLED"] = "0"

    ldflags = f"-X main.version={version}"

    print(f"🔨 Compilando {binary_name}...")

    result = subprocess.run(
        ["go", "build", "-ldflags", ldflags, "-o", str(output_path), "."],
        cwd=PROJECT_ROOT,
        env=env,
        capture_output=True,
        text=True
    )

    if result.returncode != 0:
        print(f"❌ Erro ao compilar {binary_name}: {result.stderr}")
        return None

    print(f"✅ Compilado: {binary_name}")
    return output_path


def build_all_binaries() -> Dict[str, Path]:
    """Compila binários para todas as plataformas"""
    results = {}
    for platform in PLATFORMS:
        path = build_binary(platform)
        if path:
            results[platform] = path
    return results


def update_release(base_url: str):
    """Atualiza metadata da release com binários atuais"""
    version = get_current_version()

    metadata = {
        "tag_name": version,
        "name": f"Release {version}",
        "published_at": datetime.utcnow().isoformat() + "Z",
        "html_url": f"{base_url}",
        "body": f"Atualização automática - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}",
        "assets": scan_binaries(base_url)
    }

    save_metadata(metadata)
    print(f"📦 Release atualizada: {version} com {len(metadata['assets'])} binários")
    return metadata


# ============ ENDPOINTS ============

@app.route("/")
def index():
    """Página inicial"""
    return jsonify({
        "service": "00cli Update Server",
        "version": get_current_version(),
        "endpoints": {
            "/latest": "Informações da última versão",
            "/download/<binary>": "Download de binário específico",
            "/health": "Health check",
            "/build": "POST - Compila e atualiza binários"
        }
    })


@app.route("/latest")
def latest():
    """Retorna informações da última release"""
    base_url = request.url_root.rstrip("/")

    metadata = load_metadata()

    # Atualizar assets com URLs corretas
    metadata["assets"] = scan_binaries(base_url)
    metadata["tag_name"] = get_current_version()

    return jsonify(metadata)


@app.route("/download/<binary_name>")
def download(binary_name: str):
    """Download de binário específico"""
    binary_path = BINARIES_DIR / binary_name

    if not binary_path.exists():
        abort(404, description=f"Binário '{binary_name}' não encontrado")

    return send_file(
        binary_path,
        as_attachment=True,
        download_name=binary_name
    )


@app.route("/health")
def health():
    """Health check"""
    return jsonify({
        "status": "ok",
        "version": get_current_version(),
        "binaries": len(list(BINARIES_DIR.glob("00cli*")))
    })


@app.route("/build", methods=["POST"])
def build():
    """Compila todos os binários e atualiza a release"""
    # Verificar token de autenticação (opcional)
    auth_token = os.getenv("UPDATE_SERVER_TOKEN")
    if auth_token:
        provided_token = request.headers.get("Authorization", "").replace("Bearer ", "")
        if provided_token != auth_token:
            abort(401, description="Token inválido")

    platforms = request.json.get("platforms") if request.is_json else None

    if platforms:
        # Compilar apenas plataformas específicas
        results = {}
        for platform in platforms:
            path = build_binary(platform)
            if path:
                results[platform] = str(path)
    else:
        # Compilar todas
        results = {p: str(path) for p, path in build_all_binaries().items()}

    base_url = request.url_root.rstrip("/")
    metadata = update_release(base_url)

    return jsonify({
        "status": "success",
        "built": list(results.keys()),
        "version": metadata["tag_name"],
        "assets_count": len(metadata["assets"])
    })


@app.route("/build/<platform>", methods=["POST"])
def build_platform(platform: str):
    """Compila binário para uma plataforma específica"""
    if platform not in PLATFORMS:
        abort(400, description=f"Plataforma '{platform}' não suportada. Use: {', '.join(PLATFORMS.keys())}")

    path = build_binary(platform)
    if not path:
        abort(500, description=f"Erro ao compilar para {platform}")

    base_url = request.url_root.rstrip("/")
    metadata = update_release(base_url)

    return jsonify({
        "status": "success",
        "platform": platform,
        "binary": path.name,
        "version": metadata["tag_name"]
    })


# ============ CLI ============

def cli_build():
    """Comando CLI para compilar binários"""
    print("🚀 Iniciando build de todos os binários...")
    results = build_all_binaries()

    if results:
        print(f"\n✅ Build concluído! {len(results)} binários gerados:")
        for platform, path in results.items():
            size = path.stat().st_size / 1024 / 1024
            print(f"   - {path.name} ({size:.2f} MB)")

        # Atualizar metadata
        metadata = update_release("http://localhost:8080")
        save_metadata(metadata)
    else:
        print("❌ Nenhum binário foi gerado")
        sys.exit(1)


def main():
    parser = argparse.ArgumentParser(description="00cli Update Server")
    parser.add_argument("--host", default="0.0.0.0", help="Host para bind (default: 0.0.0.0)")
    parser.add_argument("--port", type=int, default=8080, help="Porta (default: 8080)")
    parser.add_argument("--debug", action="store_true", help="Modo debug")
    parser.add_argument("--build", action="store_true", help="Apenas compilar binários (não inicia servidor)")

    args = parser.parse_args()

    if args.build:
        cli_build()
        return

    print(f"""
╔══════════════════════════════════════════════════════════════╗
║           00cli Update Server                                 ║
╠══════════════════════════════════════════════════════════════╣
║  Endpoints:                                                   ║
║    GET  /latest              - Info da última versão          ║
║    GET  /download/<binary>   - Download de binário            ║
║    POST /build               - Compilar todos os binários     ║
║    POST /build/<platform>    - Compilar plataforma específica ║
║    GET  /health              - Health check                   ║
╚══════════════════════════════════════════════════════════════╝
    """)

    print(f"🚀 Servidor iniciando em http://{args.host}:{args.port}")
    print(f"📁 Binários em: {BINARIES_DIR}")
    print(f"📦 Versão atual: {get_current_version()}")
    print()

    app.run(host=args.host, port=args.port, debug=args.debug)


if __name__ == "__main__":
    main()
