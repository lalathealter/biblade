# biblade v1.0.0 (?)

***Biblade*** (*Bible + Blade*) is a desktop utility app that allows you to quickly 
paste any of your favourite quotes - for those (presumably, *gaming*) moments when 
you literally have no time to google or to type.

This app operates directly in your current text prompt, so that you do not have
to `ALT-TAB` to employ it. Biblade is meant to serve you both like a sword of the Holy Spirit, 
which is the Word of God (Eph. 6:17), a discerner of the thoughts and intents of the heart (Heb. 4:12),
and like a shiv that you keep in your sleeve for *certain* conversations, hence its' name.

The app has been lazily tested on Windows. Primary functionality seems to be working just fine.

## Running instructions

To use biblade you can either download a github realise or build it from source.
If you prefer the latter, please clone this repository (while making sure to
have golang compiler installed).

Run:
```
go run .
```

Build an executable:
```
go build .
```

Leave this app running in the background while you intend to use it and simply
close the process (`Ctrl-C`) when it is no longer needed.

## On using biblade

While an app is running, type a predetermined key (the default is `]`) to enter into
active mode. Doing so displays an in-prompt interface (please make sure to use while
having a cursor set on text prompt) for you to choose a quote you want to paste. 
All phrases are broken into "wheel" frames — structures that hold a certain
amount (not bigger than it is set to be) of actual quotes and pointers to other
frames (which are meant to represent subsets or subthemes — to make your navigation 
easier). All items inside a wheel frame have a corresponding character displayed
in square brackets before its' tag, by typing which you choose them over others.
Choosing a quote pastes its' full text right into the same text prompt and exit
the active mode, while choosing a "next" option or a frame pointer will leave
active mode intact until you reach the end.
Typing anything aside from suggested keys (except backspace) will switch mode 
back to inactive. 

## On adding your phrases

To write your own phrase set, please open your preferred text editor, create an 
empty json file and use the following structure:

A file must start with an array of pairs. A pair is a json array of length 2 that
contains either a complete phrase or a subset of phrases. The first element of a
pair must be a string that reflects its shortened name (tag), while the second
should be either a string (a full text of a phrase to be pasted) or an array of
pairs (a subset). As you can imagine, it is possible to nest them however you like or find
best suitable.

An example:

```json
[ 
    ["your_tag_name", "full_text_of_a_phrase"],
    ["short_name", "another_quote"],
    ["#collection_name", [
        ["#1quote", "blah-blah-blah"],
        ["#2quote", "qwerty"],
        ["Animal_sounds", [
            ["cow", "moo-moo"],
            ["cat", "meow, meow"],
            ["dog", "woof"]
        ]]
    ]],
    ["cooking_recipes", [
        ["omelet", "eggs, milk..."],
        ["pasta", "spaghetti"]
    ]]
]
```

## On settings

You can change app's settings by simply editing the corresponding json file (`settings.json`) 
in your preferred text editor. Please consider the original types of values before changing them.
After you are done, please make sure to reload the app. Here is a quick rundown on what every field means:

* `frameCap` (type *int*) — the maximum allowed length of a frame; If a phrase set is longer, 
it is going to be broken into several frames. It cannot be set lower than 4
* `introLen` (type *int*) — the maximum allowed length of a tag; all longer tags are going to 
be truncated with dots. It cannot be set lower than 4.
* `wheelFile` (type *string*) — the path to a json file containing a collection your selected phrases. 
* `key` (type *string*) — the first character of a string you put in this field is going to 
be used as an activation key.


