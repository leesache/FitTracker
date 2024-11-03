# Fitness Tracker

## Overview

The Fitness Tracker is a Go-based application designed to help users track their workout routines. It reads exercise data from an Excel file and organizes it into a structured format for easy management and visualization. The application is being developed with plans for a user interface (UI) to enhance user experience.

## Features

- **Excel Import**: Import workout data from `.xlsx` files to streamline the data entry process.
- **Structured Data**: Organize training data into structured formats, making it easy to access and manipulate.
- **User Interface**: A graphical user interface (UI) will be developed to facilitate user interaction.
- **Cross-Platform Compatibility**: Designed to work seamlessly on macOS, ensuring accessibility for a wide range of users.

## Requirements

- **Go**: Ensure that Go is installed on your system. You can download it from [golang.org](https://golang.org/dl/).
- **Excelize**: The project uses the `excelize` library for reading Excel files. It can be installed using:
  ```bash
  go get github.com/xuri/excelize/v2
  ```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/fitness-tracker.git
   ```
2. Navigate to the project directory:
   ```bash
   cd fitness-tracker
   ```
3. Install the necessary Go packages:
   ```bash
   go mod tidy
   ```

## Usage

1. Prepare your Excel file with the required workout data.
2. Update the `excelImportData` function in the code to point to your specific Excel file.
3. Run the application:
   ```bash
   go run main.go
   ```

## Future Work

- Development of the graphical user interface (UI) to improve user interaction.
- Enhancements to data visualization and reporting features.
- Expansion of the application to support additional fitness tracking features.

## Contributing

Contributions are welcome! If you would like to contribute, please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For questions or inquiries, please contact [your email or contact information].
