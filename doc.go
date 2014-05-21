/*
captainhook is a generic webhook listener

This tool was built as part of a CI orchestration process, to be called when
Docker trusted builds finish.  It explicitly ignores the posted data from the webhook
because that would be `insecure`, which is `bad`.

Despite our intended purpose, it can be used to trigger any process when you receive a post
to a specific URL.  That's why we called it a generic webhook listener.

To use captainhook, first create a directory to store the json scripts that describe your
orchestration.  We'll refer to that directory as `configdir`.

Each script you create in the `configdir` will be executed when
the corresponding endpoint is called.

`mkdir ~/captainhook`

Now add a json file in the `configdir`.  There's a sample in the `example` directory with
some sterile commands that won't modify your filesystem.

Run captainhook with a command similar to this:

`captainhook -configdir ~/captainhook`

If you have a script called `deployBigApp.json` you would trigger
it by posting to http://your.captainhook.url/deployBigApp.

The scripts in the json file are executed sequentially, and the output is logged
and returned to the caller in the response, which always has an HTTP status code
of 200 (OK) even if your scripts didn't work.  This is intentional, to avoid causing
errors in external services like Docker or Github, which might not like you returning
statuses other than 200 (OK).

Copyright 2014, Brian Ketelsen and Kelsey Hightower

LICENSE information found in LICENSE file.

*/
package main
