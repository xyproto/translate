# translate

Transate text from `stdin` to the current locale using `ollama` and `mixtral:instruct`.

`mixtral:instruct` (around 26 GiB) will be downloaded the first time the utility is being run.

Requires `ollama` to be set up and serving on port `11434`.

The translation takes a cuple of seconds, and is not terribly exact, but it is a pretty versatile tool.

`ollama` and `translate` does not require an internet connection once `mixtral:instruct` has been downloaded.

## Example use

`LANG` is set to `nb_NO.UTF-8`

```
$ echo 'I can speak Norwegian!' | ./translate
Jeg kan snakke norsk!
```

`LANG` is explicitly set to `de_DE`:

```
$ echo 'I can speak German!' | LANG=de_DE ./translate
Ich kann Deutsch sprechen!
```

`LANG` is explicitly set to `de_DE` and verbose output is enabled with `-v`:

```
$ echo 'I can speak German!' | LANG=de_DE ./translate -v
Prompt: Translate the following text to the locale de_DE (and only output the translated text): I can speak German!
Sending request to /api/tags
Sending request to /api/generate: {"model":"mixtral:instruct","prompt":"Translate the following text to the locale de_DE (and only output the translated text): I can speak German!"}
Ich kann Deutsch sprechen!
```

## General info

* Version: 1.0.0
* License: BSD-3
