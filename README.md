Golang project to create terraform helper utility

Terraform is an awesome way to create infrastructure. The DSL is
very easy to pick up, and there are abundant examples, as well
as a decent number of modules.


I wanted to simplify it ever so much for operators and not
 _require_ them to be able to read and write terrafrom code.
I also put in pieces that are my opinionated version of best
practices.  


More specifically, I very much like [OpenCredo's Terraform Design Patterns](https://opencredo.com/terraform-infrastructure-design-patterns/)
for modules, but tried to take it a bit further.

For all languages modules have is a way to list out the
set of required inputs and the outputs for a given function/operation/classes

Usually those langauges also have tools for generating documentation of those inputs and
outputs over and above the presumably great helpdocs also written.

You an group those together, like all of the public classes in a maven module,
or ruby gem, etc.  Most of these have their own package manager that will grab remote
modules and install them locally for your use.  This is done through a manifest
(e.g. pom.xml, Gemspec, requirements.txt, package.json).  those
languages also usually have tools so that you can see the various function/operation/classes
and know what variables are required, what optional values are (maybe), and what
the expected output is for those.

Finally their are the build tools and conventions, i.e. maven, gradle, bundler,
npm, etc.  

This tool chain is meant to start filling in some of those gaps.  In all likelihood
the Hashicorp team will have better answer to all of these issues, and this tool
will be obviated.  In the mean time however, I got a bit more productive at using
terraform, and finally did a real project with golang.

It has a basic package retriever,  
inspector for modules that tell the inputs and Outputs, as well as boilerplate
generators for common patterns.


groups of those things, like all the public classes in a package in java,
or a ruby gem library that contains many of those definitions.


Projects -> Layers -> Modules

modules have clients

Module Invoke




fetch module to local Copyright
Create Documentation - yaml doc
doc points at original for info, sha, and local location.



Same doc -

[See Usage](docs/formterra.md)

./formterra module inspect -u tfproject/test-fixtures/simpleterraform -u git::ssh://git@github.com/jmahowald/terraform-aws-vpc.git//modules/network --name joshtest | ./formterra module client
