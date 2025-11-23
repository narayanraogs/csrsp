# CSRSP (Checkout Software for Remote Sensing Payloads)

## Overview
CSRSP is a comprehensive satellite ground station software designed for the monitoring, control, and data processing of remote sensing payloads. It features a robust **Go** backend for data handling and a modern **Flutter** web frontend for user interaction.

## Key Features

### 1. Data Acquisition
- **Modes**: Support for Frame-based, Time-based, and Continuous acquisition.
- **Real-time Monitoring**: Live updates on frame counts and data integrity.
- **Hardware Integration**: Interfaces with satellite payloads and ground station hardware.

### 2. Offline Processing
- **Data Management**: Browse, filter, and manage historical acquisition data.
- **Reprocessing**: Apply different Result Profiles to previously acquired raw data.
- **Reporting**: Generate and view detailed intermediate and final reports.

### 3. Bit Error Rate (BER) Logging
- **Real-time Plotting**: Visualize up to 6 parameters across 6 concurrent streams.
- **Long-duration Testing**: Supports logging sessions up to 2 hours with high-frequency sampling.
- **Archival**: Automatic storage of test data for analysis.

### 4. Administration & Security
- **Role-Based Access**: Granular permissions for Acquisition, Processing, and Admin tasks.
- **Dual Authentication**: Support for IP Whitelisting and User Credential login.
- **System Monitoring**: Dashboard for server health (CPU, Memory) and connection status.

## Architecture

### System Overview
- **Frontend**: Flutter Web application (`client/`) providing a responsive and interactive dashboard.
- **Backend**: Go server (`server/`) serving the web app and hosting the data processing engine.
- **Deployment**: Single binary deployment (Flutter assets embedded in Go binary).

### Data Processing Pipeline (`server/processor`)
The core of CSRSP is a high-performance, modular data pipeline designed for satellite telemetry:
- **Stage-Based Design**: Data flows through a series of independent stages (e.g., Frame Synchronization, FEC, Parameter Extraction).
- **Standards Compliance**:
    - **CCSDS**: Built-in support for CADU extraction, Virtual Channel routing, and Reed-Solomon error correction.
    - **MIL-STD-1750A**: Native decoding of legacy floating-point formats.
- **Extensibility**: The `GenericProcessor` stage allows offloading complex algorithms to external executables via Unix domain sockets, enabling polyglot processing modules.
- **Validation Engine**: A robust parameter validation system supporting Analog, Radix, CRC, and Status checks with real-time fan-out processing.

## Getting Started

### Prerequisites
- **Go**: Version 1.18 or higher.
- **Flutter**: Latest stable channel.
- **Make**: For running build commands.

### Building the Project
The project uses a `Makefile` to automate the build process.

1.  **Build Everything (Client + Server)**
    ```bash
    make all
    ```
    This command will:
    - Compile the Flutter web application.
    - Copy the assets to the server's web directory.
    - Compile the Go server with version stamping.

2.  **Build Client Only**
    ```bash
    make client-build
    ```

3.  **Build Server Only**
    ```bash
    make server-build
    ```

4.  **Clean Build Artifacts**
    ```bash
    make clean
    ```

## Project Structure
- `client/`: Flutter web application source code.
- `server/`: Go server source code.
  - `processor/`: Telemetry processing pipeline and stages.
  - `utils/`: Shared utilities (Bit manipulation, MIL-STD-1750A, Stats).
  - `db/`: Database access layer and SQLC generated code.
- `workflows/`: Functional specifications and workflow documentation.
- `Makefile`: Build automation scripts.
