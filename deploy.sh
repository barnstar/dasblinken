#!/usr/bin/env bash
set -euo pipefail

# Defaults (override via env vars or flags)
HOST="${HOST:-minipi}"
PIUSER="${PIUSER:-admin}"
DEST_DIR="${DEST_DIR:-/home/admin/piled}"
CONFIG="${CONFIG:-config.json}"
AUTHKEY="${AUTHKEY:-authkey.txt}"

usage() {
  cat <<EOF
Usage: $(basename "$0") <command> [options]

Commands:
  bin             Copy ./piled to remote 
  cfg             Copy ./${CONFIG} to remote 
  effects         Copy ./effects/effects.json to remote
  authkey         Copy ./${AUTHKEY} to remote DEST_DIR/${AUTHKEY}
  all             Run bin, cfg, effects, authkey
  run             Stop existing piled and start it with --config=${CONFIG}
  stop            Stop piled on remote
  ssh             Open an ssh shell to remote
  show           Show remote config and effects.json contents
  which           Show current resolved settings

Options (env vars or flags):
  --host NAME        Remote host (default: $HOST)
  --piuser NAME      Remote PIUSER (default: $PIUSER)
  --dest PATH        Remote dest dir (default: $DEST_DIR)
  --config FILE      Config file (default: $CONFIG)
  --authkey FILE     Authkey file (default: $AUTHKEY)

Examples:
  HOST=raspberrypi PIUSER=pi ./deploy.sh all
  ./deploy.sh --host minipi --piuser admin run
EOF
}

# Parse flags
while [[ $# -gt 0 ]]; do
  case "$1" in
    --host) HOST="$2"; shift 2 ;;
    --piuser) PIUSER="$2"; shift 2 ;;
    --dest) DEST_DIR="$2"; shift 2 ;;
    --config) CONFIG="$2"; shift 2 ;;
    --authkey) AUTHKEY="$2"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    bin|cfg|effects|authkey|all|run|stop|ssh|show|which)
      CMD="$1"; shift; break ;;
    *) echo "Unknown arg: $1"; usage; exit 1 ;;
  esac
done

REMOTE="${PIUSER}@${HOST}"

ensure_remote_dir() {
  ssh "$REMOTE" "mkdir -p '$DEST_DIR'"
}

cmd_bin() {
  if [[ ! -f ./piled ]]; then
    echo "Error: ./piled not found. Build first: make piled"
    exit 1
  fi
  ensure_remote_dir
  rsync -ave ssh ./piled "$REMOTE:$DEST_DIR/piled"
}

cmd_configs() {
  if [[ ! -f "./${CONFIG}" ]]; then
    echo "Error: ${CONFIG} not found in repo root."
    exit 1
  fi
  ensure_remote_dir
  rsync -ave ssh "./${CONFIG}" "$REMOTE:$DEST_DIR/${CONFIG}"
}

cmd_effects() {
  if [[ ! -f ./effects/effects.json ]]; then
    echo "Error: ./effects/effects.json not found."
    exit 1
  fi
  ensure_remote_dir
  rsync -ave ssh ./effects/effects.json "$REMOTE:$DEST_DIR/effects.json"
}

cmd_authkey() {
  if [[ ! -f "./${AUTHKEY}" ]]; then
    echo "Warning: ${AUTHKEY} not found; skipping."
    return 0
  fi
  ensure_remote_dir
  rsync -ave ssh "./${AUTHKEY}" "$REMOTE:$DEST_DIR/${AUTHKEY}"
}

cmd_all() {
  cmd_bin
  cmd_configs
  cmd_effects
  cmd_authkey
}

cmd_run() {
  ssh -n -f "$REMOTE" "sh -c 'sudo killall -9 piled > /dev/null 2>&1 || true; cd \"$DEST_DIR\"; nohup sudo ./piled --config=\"$CONFIG\" > /dev/null 2>&1 &'"
  echo "Started piled on $REMOTE with --config=$CONFIG"
}

cmd_stop() {
  ssh "$REMOTE" "sudo killall -9 piled > /dev/null 2>&1 || true"
  echo "Stopped piled on $REMOTE"
}

cmd_ssh() {
  exec ssh "$REMOTE"
}

cmd_show() {
  echo "Remote: $REMOTE"
  echo "Directory: $DEST_DIR"
  echo "Config file: $CONFIG"
  echo
  echo "=== ${DEST_DIR}/${CONFIG} ==="
  ssh "$REMOTE" "test -f '$DEST_DIR/$CONFIG' && cat '$DEST_DIR/$CONFIG' || echo 'No config found.'"
  echo
  echo "=== ${DEST_DIR}/effects.json ==="
  ssh "$REMOTE" "test -f '$DEST_DIR/effects.json' && cat '$DEST_DIR/effects.json' || echo 'No effects.json found.'"
}

cmd_which() {
  cat <<EOF
HOST=$HOST
PIUSER=$PIUSER
DEST_DIR=$DEST_DIR
CONFIG=$CONFIG
AUTHKEY=$AUTHKEY
REMOTE=$REMOTE
EOF
}

case "${CMD:-}" in
  bin)            cmd_bin ;;
  cfg)            cmd_configs ;;
  effects)        cmd_effects ;;
  authkey)        cmd_authkey ;;
  all)            cmd_all ;;
  run)            cmd_run ;;
  stop)           cmd_stop ;;
  ssh)            cmd_ssh ;;
  which)          cmd_which ;;
  show)           cmd_show ;;
  *)              usage; exit 1 ;;
esac