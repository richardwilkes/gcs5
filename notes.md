### JSON encoding/decoding

The JSON encoding/decoding in this repo is done essentially manually (i.e. not making use of Go's standard json tagging
mechanism). This is due to the legacy Java code and how it conditionally loaded and stored data as well as how some
objects relied on other information during both load and store, which is not supported by Go's standard implementation.

### Field access

When the code was in Java, all fields were accessed through getters and setters. Those have been eliminated where the
field itself doesn't have special requirements when being read or written, which greatly reduces a lot of boiler-plate
accessor code.
