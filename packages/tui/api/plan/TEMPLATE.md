
---

**Goal**: Write a Python script that scrapes a given URL and extracts all image links.
**Sub-goals**:
1.  Fetch the HTML content of the target URL.
2.  Parse the HTML to find all image tags (`<img>`).
3.  Extract the `src` attribute from each image tag.
4.  Return the list of image links.

## Specifications

*   **Input**: A single string variable, `url`, containing the target URL.
*   **Output**: A Python list of strings.
*   **Dependencies**: The script should use the `requests` and `beautifulsoup4` libraries.
*   **Error Handling**: If the URL is invalid or the request fails, the script should raise a descriptive exception.

## Context

**Tools**: Python interpreter, file system access, requests library, beautifulsoup4 library.

## Constraints

*   **Security**: Do not fetch resources from non-HTTP/HTTPS protocols.
*   **Performance**: Optimize for speed on single-page requests.
*   **Code Style**: Adhere to PEP 8 standards.

## Plan Format

1.  **Task**: Fetch URL Content
    *   **Steps**: Use the `requests` library to perform a GET request on the input URL. Check the response status code.
2.  **Task**: Parse HTML
    *   **Steps**: Initialize a `BeautifulSoup` object with the response content.
3.  **Task**: Extract Image Links
    *   **Steps**: Use `BeautifulSoup` to find all `img` tags. Loop through the tags and get the `src` attribute.

## Execution

1.  **Workflow**: Follow the plan sequentially.
2.  **Validation**: After each task, validate the output against the specifications. For example, check that the "Extract Image Links" task returns a list of strings.
3.  **Reporting**: Log progress and final results.
4.  **Refinement**: If the plan fails, review the logs and attempt to generate an alternative plan to solve the problem.

