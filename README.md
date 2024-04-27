# AutoEx

AutoEx is an open-source utility designed to monitor a specified directory for compressed files and automatically attempt to extract them using a list of provided passwords. It supports recursive extraction, ensuring that all nested compressed files within the directory are also processed until there are no more compressed packages left to extract.

## Features

- **Automatic Monitoring:** Continuously watches a specified folder for new compressed files.
- **Password List Support:** Tries multiple passwords from a provided list to extract password-protected archives.
- **Recursive Extraction:** Automatically extracts nested compressed files.
- **7z Volume Extraction Support:** Capable of extracting multi-volume 7z archives.

## Getting Started

### Prerequisites

- Golang
- Docker (for container-based deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ChenyuanHu/autoex.git
   ```
2. Navigate to the cloned directory:
   ```bash
   cd autoex
   go build
   ```

### Docker Setup

#### Building the Docker Container

To containerize AutoEx, build the Docker image using:

```bash
docker build -t autoex .
```

This will compile the AutoEx application within a Docker container, ensuring that it runs in an isolated and consistent environment.

#### Running AutoEx in a Docker Container

After building the image, run AutoEx using:

```bash
docker run -d \
  -v /path/to/your/data:/data \
  -e AUTOEX_DIR="/data" \
  -e AUTOEX_PW_LIST="password1|password2|password3" \
  -e AUTOEX_DEL_COMPLETE="false" \
  --name autoex_container autoex
```

This command sets up AutoEx to monitor the `/data` directory inside the container, which maps to `/path/to/your/data` on your host machine. It will attempt to extract any new archives using the provided password list.

## Configuration

AutoEx is configured entirely through environment variables, allowing you to control its behavior without modifying the code. Here are the environment variables you can set:

- **`AUTOEX_DIR`**: Specifies the directory that AutoEx will monitor for compressed files. For example, to watch the folder `/home/user/downloads`, you would set it like this in a Unix-like system:
  ```bash
  export AUTOEX_DIR="/home/user/downloads"
  ```

- **`AUTOEX_PW_LIST`**: Sets the list of passwords that AutoEx will try when attempting to extract protected archives. Passwords should be separated by the `|` character. Include an empty password in the list if you want AutoEx to attempt to extract archives that might not be password protected:
  ```bash
  export AUTOEX_PW_LIST="password1|password2|password3"
  ```

Ensure these environment variables are set before running AutoEx. You can place these export commands in your shell's startup file (like `.bashrc` or `.zshrc`) to make the configuration persistent across sessions, or define them temporarily just before running the script.

### Usage

To start using AutoEx, you need to specify the directory to monitor and the path to your password list file. Run the tool with the following command:

```bash
  AUTOEX_DIR="/home/user/downloads" AUTOEX_PW_LIST="password1|password2|password3" AUTOEX_DEL_COMPLETE="false" ./autoex
```

Or (not be password protected)
```bash
  AUTOEX_DIR="/home/user/downloads" AUTOEX_PW_LIST="" AUTOEX_DEL_COMPLETE="false" ./autoex
```

## Contributing

Contributions to AutoEx are welcome! Here's how you can contribute:

- **Fork the repository:** Create your own copy of the project.
- **Create a new branch:** `git checkout -b new-feature`
- **Make your changes:** Add your new feature or improvement.
- **Submit a pull request:** Open a PR against the main branch, and provide a description of your changes.

Please ensure your code adheres to the existing style to maintain consistency.

## License

AutoEx is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Support

If you encounter any bugs or have suggestions for improvements, please file an issue on the GitHub issues page.

Thank you for using or contributing to AutoEx!