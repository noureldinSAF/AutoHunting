import os
import yaml

CONFIG_PATH = "config/config.yaml"


def load_config():
    if os.path.exists(CONFIG_PATH):
        with open(CONFIG_PATH, "r") as file:
            return yaml.safe_load(file) or {}
    return {}


def ensure_default_config():
    os.makedirs("config", exist_ok=True)
    if os.path.exists(CONFIG_PATH):
        return
    default_config = {"server": {"host": "0.0.0.0", "port": 8000}}
    with open(CONFIG_PATH, "w") as file:
        yaml.dump(default_config, file, default_flow_style=False)
    print(f"Created default config file at {CONFIG_PATH}")
