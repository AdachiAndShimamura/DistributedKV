import os
import pathlib
import subprocess
import sys
from pathlib import Path


def get_root_directory():
    """
    获取当前脚本的根目录，即脚本所在目录的父目录的父目录。
    """
    # 获取当前脚本的绝对路径
    script_path = Path(__file__).resolve()
    # 获取父目录的父目录
    root = script_path.parent.parent
    return root


def execute_command(command):
    """
    执行指定的命令，并处理可能的错误。
    
    :param command: 要执行的命令字符串
    """
    print(f"正在执行命令: {command}")
    try:
        # 使用 subprocess.run 执行命令
        result = subprocess.run(
            command,
            shell=True,  # 使用 shell 执行
            check=True,  # 如果命令返回非零退出码，则引发 CalledProcessError
            stdout=subprocess.PIPE,  # 捕获标准输出
            stderr=subprocess.PIPE,  # 捕获标准错误
            text=True  # 以字符串形式返回输出
        )
        # 打印标准输出
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        # 打印错误信息
        print(f"执行命令时出错: {command}", file=sys.stderr)
        print(e.stderr, file=sys.stderr)
        sys.exit(e.returncode)


def main():
    root = get_root_directory()
    print(f"根目录: {root}")
    os.chdir(root)

    raw = (
        "protoc -I=proto/proto "
        "--go_opt=paths=source_relative --go_out=proto/gen "
        "--go-grpc_opt=paths=source_relative --go-grpc_out=proto/gen "
        "./proto/proto/kvrpcpb/*.proto"
    )

    raft = (
        "protoc -I=proto/proto "
        "--go_opt=paths=source_relative --go_out=proto/gen "
        "--go-grpc_opt=paths=source_relative --go-grpc_out=proto/gen "
        "./proto/proto/raftpb/*.proto"
    )

    server = (
        "protoc -I=proto/proto "
        "--go_opt=paths=source_relative --go_out=proto/gen "
        "--go-grpc_opt=paths=source_relative --go-grpc_out=proto/gen "
        "./proto/proto/*.proto"
    )
    commands = [raw, raft, server]

    if not pathlib.Path(root/"proto"/"gen").exists():
        os.mkdir(root/"proto"/"gen")


    for cmd in commands:
        execute_command(cmd)


if __name__ == "__main__":
    main()
