# Getting support

I am willing to provide limited support for this project using this project's issue tracker. Keep in mind that I maintain this project in my spare time.

Before opening an issue, please do the following first:

1. Check the [Frequently Asked Questions][faq] in the README file. These should cover the most common issues that you might encounter.
1. If you encountered an error message, **please read that error message thouroughly** and apply common sense.
1. Make sure your configuration is correct. Pay special attention to the following points:
    1. Are all file paths (access log files) configured correctly?
    1. Is your access log format configured correctly?
    1. Are required file system permissions set correctly?
1. Check other open issues for similar requests.

Should you not be able to resolve your issue on your own using the steps above, please open a new issue with the ![label: question][~question] label.

When opening a new issue, please provide the following information:

- A clear and concise description of the expected and actual behaviour of the program
- The contents of your configuration file (or command-line flags, if you're not using a config file)
- The STDOUT and STDERR output of the exporter process
- An example log file. **Please limit the length to as few lines necessary to reproduce your issue**. Also, please make sure to redact any sensitive data (especially [PII][pii]) from your log files (IP addresses, user names, email addresses, etc.).
- The version of the exporter that you're using

I'll try to react to support requests in a timely manner. If your issue should not receive any attention after a reasonable amount of time, feel free to ping me in that issue.

[faq]: https://github.com/martin-helmich/prometheus-nginxlog-exporter#frequently-asked-questions
[pii]: https://en.wikipedia.org/w/index.php?title=Personally_identifiable_information&redirect=no
[~question]: https://img.shields.io/badge/-question-cc317c.svg