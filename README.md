# oldhosts

`oldhosts` is a tool for bug bounty hunters to discover old hosts that are no longer available, but might still be present on different known and related servers.

## Installation

To install `oldhosts`, follow the steps below:

    - Ensure you have Go installed on your system.
    - Run the following command to install the required packages:
    ```
    go install -v github.com/topscoder/oldhosts
    ```

## Options
Run the script using the following command-line arguments:

    ```
    oldhosts -ips <ips> -hosts <hosts> [-curl] [-silent]
    ```

    - `-ips` (required): Specify an IP address or provide a filename containing IP addresses (one per line).
    - `-hosts` (required): Specify the hostname or provide a filename containing hosts (one per line).
    - `-curl` (optional): Output the results as Curl commands.
    - `-silent` (optional): Run in silent mode, suppresses non-200 responses (except for content length).

View the results:

    - The script will perform HTTP and HTTPS requests for each IP and host combination.
    - The script will display the response status code and content length for each successful request.
    - If the `-curl` flag is specified, Curl commands will be displayed instead of the response details.

## Example

Here is an example command to run `oldhosts`:

```
oldhosts -ips "192.168.0.1" -hosts "example.com" -curl
```

This command will perform HTTP and HTTPS requests to the specified IP addresses and hosts, displaying the results as Curl commands.

## Notes

- The script limits the number of concurrent calls to 5 for performance reasons. You can adjust this value by modifying the `semaphore` channel in the code.

- The script supports both individual strings and filenames as input for IP addresses and hosts. If a filename is provided, the script reads the IP addresses and hosts from the file (one per line).

- The script removes any trailing slashes from the hosts and tries to append default ports (":80" for HTTP and ":443" for HTTPS) to the host header.

- The script has a timeout of 1 second for each HTTP request.

- If the `-silent` flag is specified, the script will only print results for successful requests (200 status code). Use this flag to reduce the output and focus on relevant information.

## Contributing

Contributions are welcome! If you find a bug or want to suggest a new feature, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
