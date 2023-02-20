# Contributing

## Go version

The project usually uses the latest Go version as declared by the `go.mod` file.
You may not be able to build it using older compilers.

## Local development

We recommend reading the `Makefile` to discover some targets which you can
execute. It can be used as a shortcut to run various useful commands.

You may have to run the following command to install a linter and a code
formatter before executing certain targets:

    $ make tools

If you want to check if the pipeline will pass for your commit it should be
enough to run the following command:

    $ make ci

It is also useful to often run just the tests during development:

    $ make test

Easily format your code with the following command:

    $ make fmt

## Naming tests

When naming tests which tests a specific behaviour it is recommended to follow a
pattern `TestNameOfType_ExpectedBehaviour`. Example:
`TestCreateHistoryStream_IfOldAndLiveAreNotSetNothingIsWrittenAndStreamIsClosed`
.

## Opening a pull request

Pull requests are verified using CI, see the previous section to find out how to
run the same checks locally. Thanks to that you won't have to push the code to
see if the pipeline passes.

It is always a good idea to try to [write a good commit message][commit-message]
and avoid bundling unrelated changes together. If your commit history is messy
and individual commits don't work by themselves it may be a good idea to squash
your changes. [Effective Go][effective-go] and [Go Code Review
Comments][code-review-comments] are good to read.

### Feature branches

When naming long-lived feature branches please follow the pattern `feature/...`.
This enables CI for that branch.


[commit-message]: https://cbea.ms/git-commit/

[effective-go]: http://golang.org/doc/effective_go.html

[code-review-comments]: https://github.com/golang/go/wiki/CodeReviewComments