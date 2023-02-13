# DOE Science Bowl PDF Parser
[![status-badge](https://ci.adawesome.tech/api/badges/ADawesomeguy/doe-scibowl-pdf-parser/status.svg)](https://ci.adawesome.tech/ADawesomeguy/doe-scibowl-pdf-parser)  
Easy-to-use program written in Go to extract DOE Science Bowl packets to raw JSON

## Scope
The DOE happens to distribute their question in PDF format with no alternative and I needed a way to get those into JSON for an API so here we are.

## Usage
You can either use the library, found in the `parse` directory for the functions regarding parsing the PDFs, or compile this yourself and spin up the web server. This web server takes form data requests to `/pdf` and expects a PDF file to which it will return a JSON object for all questions.

## Licensing and Credits
[@ADawesomeguy](https://github.com/ADawesomeguy) wrote this and is a very cool guy who happened to license it under the [MIT](https://opensource.org/licenses/MIT) open-source license.
