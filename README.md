`envare` - safe `.env` file parsing to use from the shell

## what is it?

`envare` is a small and simple golang program to read `.env` files and emit
some or all of their contents to stdout, where they can be easily consumed by
shell scripts.

## why is it for?

I often want to read from `.env` files from shell scripts.

I know I could just say `. .env` but I've never felt comfortable doing so.
Even if I control those `.env` files, and I usually do, `source`ing them has
always seemed like asking for trouble.  So, in practice, I haven't done that.
Instead I end up hardcoding values in my convenience scripts and I don't
commit them to version control.

Not ideal.

Here's a little script that I use nearly every day, often from inside
containers, to launch `psql`:

```shell
#!/bin/bash

echo 'db_server:5432:db_name:db_user:db_password' >> $HOME/.pgpass
chmod 600 $HOME/.pgpass

psql -P expanded=auto -P null=NULL -U db_server -h db_server db_name "$@"
```

It saves me some typing, but it's still far from ideal.  And in an unwary
moment I might commit it!

With `envare`, though, I can make it generic and less hacky.

```shell
#!/bin/bash

export PGPASSWORD="$(envare DB_PASSWORD)"

psql -P expanded=auto -P null=NULL \
    -U "$(envare DB_USER) \
    -h "$(envare DB_SERVER)" \
    "$(envare DB_DATABASE) \
    "$@"
```

If launching all those subshells bothers you, first: it shouldn't; and second:
I mean it, it shouldn't; and third: you could run `envare` just once without
arguments and capture-and-eval its output.


## why shouldn't I use it?

- It's brand new.  And even after I've used it myself for a while I don't
  think I'll call it battletested because I hate that term.
- There are no tests.
- It uses [godotenv](https://github.com/joho/godotenv) and
  [shellescape](https://github.com/alessio/shellescape) to do nearly all the
  heavy lifting and maybe you don't trust them or like how they work for some
  reason.
- Maybe it's missing a feature you could use.  (In particular, it'd be nice to
  teach it to use an arbitrary file instead of whatever's in `$PWD`.)

