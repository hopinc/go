# Contributing to hop-go

This repository welcomes contributors! There are a few repository specific things to be aware of before contributing.

## Generating categories
Adding the actual dot separated categories is quite annoying to do by hand in Go, so the functionality to do this is auto-generated. If you wish to add a new category, you should open the file `categories.json` and add it in the format `"BaseCategory": ["SubCategories"]`.

After you run `make generate`, you will notice the categories have been generated in the format `ClientCategory<BaseCategory>[SubCategory]` and there is a new `new<BaseCategory>` function.

If you are adding a base category, you will want to go into `client.go` and add `<BaseCategory> *ClientCategory<BaseCategory>` to the `Client` struct. After this, in the new Client function, you will want to call the new function (`<BaseCategory>: new<BaseCategory>(&c)`).

## Adding tests to the project
If you add a new function, you will likely desire to test it. Please place the tests into the `<filename>_test.go` file. When you write your tests, it is preferred where possible you use the test tables format. You can then run `make test` to run the test suite. If you wish to view your test coverage, you can run `make cov-html`. This will open a browser to a page where you can navigate around and view the test coverage.

## Adding/updating JSON validation tests to types
If you make/edit a type in the `types` package, you will probably want to add the structure to the test suite and update it. To do so, simply open `types/types_test.go` and find the comment for the file specified. If the file isn't present, make `// <filename>.go` and then add under that in the slice (`reflect.TypeOf(<value>)`). From here, simply run `make update-types` and the JSON types will be added/updated.
