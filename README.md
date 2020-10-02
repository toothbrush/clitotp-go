# clitotp-go

This is functionally meant to be a clone of
https://github.com/toothbrush/clitotp, except now it's in Golang, so
it's not a pain for me to install in NixOS.  I'll put the original
README here for your perusal.

This is a very hacky tool to generate TOTPs on the CLI in a way that
i understand and can manage.  It puts the data where i want it,
encrypted with GnuPG.  This will probably not be suitable for your
needs, but it scratches my itch.  I encourage you to look at one of
the many other superb CLI TOTP utilities out there if this is too
rough and ready for you.

# Instructions for the impatient

You need Golang.  This worked for me:

```ShellSession
$ cd src/clitotp-go
$ go install
```

It'll now be in your `$GOPATH/bin`, so remember to add that to your
`$PATH`.

Have a look at `clitotp-go --help` for more information.

## Using the tool

You'll want to start by inserting a secret into the database.  The
file locations are hard-coded, deal with it.  The only argument these
scripts take (for now) is the name of the TOTP thing.  Probably the
website or app name is a good choice here.

```ShellSession
$ clitotp-go add github
Will insert into: /Users/yourfineface/.totp/github.gpg
Give me the secret (C-c cancels): aaaa bbbb cccc dddd
Encrypted and saved.
```

When it comes time for a site to nag you about 2FA, you can call up
the relevant thing like so:

```ShellSession
$ clitotp-go generate github
Trying to decrypt /Users/yourfineface/.totp/github.gpg... done.
164160
```

# Appendix: completion

I am a lazy person, so can't be bothered remembering or typing the
names of my TOTP secrets.  Thankfully the excellent [`cobra`
library](https://github.com/spf13/cobra/) makes completions easy.
This work for me.  It may help you, too.

```ShellSession
$ clitotp-go completion zsh > "${fpath[1]}/_clitotp-go"
```

Peace, love, and vegetables, or something like that. 🌽
