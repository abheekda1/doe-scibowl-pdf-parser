# DOE Science Bowl PDF Parser
![CI Status Badge](https://ci.adawesome.tech/api/badges/ADawesomeguy/doe-scibowl-pdf-parser/status.svg)  
Easy-to-use program written in Go to extract DOE Science Bowl packets to raw JSON

## Scope
The DOE happens to distribute their question in PDF format with no alternative and I needed a way to get those into JSON for an API so here we are.

## Usage
You can either use the library, found in the `parse` directory for the functions regarding parsing the PDFs, or compile this yourself and spinning up the web server. This web server takes form data requests and expects a PDF file to which it will return a JSON object.

## Licensing and Credits
[@ADawesomeguy](https://github.com/ADawesomeguy) wrote this and is a very cool guy who happened to license it under the [MIT](https://opensource.org/licenses/MIT) open-source license.