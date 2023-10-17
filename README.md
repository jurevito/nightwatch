# Google Search Results Parser
`main.go` is a script that extracts hyperlink data from Google search result HTML. The repository includes an example input file `pizza.html` that has 100 results for the search word _**pizza**_, example configuration file `config.json` that defines CSS selectors for links we want to parse, and example output file `output.json` that is produced when running the script.

## Usage
You have to change configuration fields so that:
- `main_elem`: selects a container element that contains all links in search results.
- `caurasel_elem`: selects a container element that contains `caurasel` or `local` links.

Other configuration fields are mostly triplets, except for `photo_img` which selects `<img>` in **Photo** caurasel. They are used like:
- `*_elem`: selects a container element that contains both the link `<a>` and element with title. Can be the same as the link itself.
- `*_link`: selects a link `<a>` that has an `href` attribute with URL of the link.
- `*_title`: selects an element whose text is the title of the co-responding link. Can be the same as the link itself.