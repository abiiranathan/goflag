package goflag

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var functions = `
check_string_in_array() {
	local string_to_check="$1"
	shift
	local array=("$@")
	local found=0
	
	for element in "${array[@]}"; do
	  if [[ "$element" == "$string_to_check" ]]; then
		found=1
		break
	  fi
	done
  
	# Return 1 if found, 0 otherwise
	return $found
}

`

// GenerateZshCompletion generates a Zsh completion script for goflag.
func (ctx *Context) GenerateZshCompletions(w io.Writer) {
	fmt.Fprintln(w, functions)

	command := filepath.Base(os.Args[0])
	fmt.Fprintf(w, "\n#compdef %s\n", command)
	fmt.Fprintf(w, "compdef _goflag_complete %s\n", command)
	fmt.Fprintln(w, "# Zsh completion script for goflag")
	fmt.Fprintln(w, "_goflag_complete() {")
	// store variable to main program
	fmt.Fprintln(w, "  local program_name=\"$words[1]\"")
	fmt.Fprintln(w, "  local words_copy=(\"$words[@]\")")
	fmt.Fprintln(w, "  words_copy=(\"${words_copy[@]:1}\")")

	fmt.Fprintln(w, "  local -a opts")
	fmt.Fprintln(w, "  local cur prev")
	fmt.Fprintln(w, "  cur=\"${words_copy[-1]}\"")
	fmt.Fprintln(w, "  prev=\"${words_copy[-2]}\"")
	fmt.Fprintf(w, "  opts=(\"help:Print help message and exit\"")

	// Generate completion options for global flags
	for index, flag := range ctx.flags {
		if index > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, `"--%s:%s"`, flag.name, flag.usage)
	}

	fmt.Fprintf(w, ")")
	fmt.Fprintln(w)

	// declare an array of valid subcommand
	fmt.Fprintf(w, "\n  declare -a valid_subcommands=(")

	// Generate completion options for subcommands and their flags
	for index, cmd := range ctx.subcommands {
		if index > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, `"%s:%s"`, cmd.name, cmd.description)
	}

	fmt.Fprintf(w, ")")
	fmt.Fprintln(w)

	for _, cmd := range ctx.subcommands {
		arrName := fmt.Sprintf("flags_%s", cmd.name)
		fmt.Fprintf(w, "  declare -a %s=(", arrName)

		// Generate completion options for subcommand flags
		for index, flag := range cmd.flags {
			if index > 0 {
				fmt.Print(" ")
			}
			fmt.Fprintf(w, "\"--%s:%s\"", flag.name, flag.usage)
		}
		fmt.Fprintf(w, ")\n")
	}

	fmt.Fprintln(w)

	script := `
  # Check if a subcommand already exists in words
  subcommand_exists=0
  current_subcmd=""
  for word in "${words_copy[@]}"; do
  	if [[ -n $word && "$word" != -* && "$word" != "--"* && "$word" != "$cur" ]]; then
  		# Found a non-flag word, assume it's a subcommand
  		subcommand_exists=1
  		current_subcmd="$word"
  		break
  	fi
  done
    
  check_string_in_array "$prev" "${valid_subcommands[@]}"
  result=$?
  
  if [[ $subcommand_exists -eq 1 || $result -eq 1 ]]; then
      local match="$prev"
      if [[ -n "$current_subcmd" ]]; then
        match="$current_subcmd"
      fi
  
      %s
      *)
        _describe 'Global Flags' opts
        ;;
      esac
      return
  fi
  
  if [[ -n "$program_name" && -z "$prev" ]]; then
  	# $program_name is not empty, and $prev is empty or null
  	local alloptions=("${valid_subcommands[@]}" "${opts[@]}")
  	_describe 'subcommands' alloptions
  else
  	# Handle other cases
  	_describe 'options' opts
  fi
  
}
`

	// Snippet for switch statement.
	var sw = new(bytes.Buffer)
	fmt.Fprintf(sw, "case \"$match\" in \n")
	for _, cmd := range ctx.subcommands {
		var arrName = fmt.Sprintf("flags_%s", cmd.name)
		fmt.Fprintf(sw, "    \"%s\")\n", cmd.name)

		fmt.Fprintf(sw, "        _describe 'Flags' %s\n", arrName)
		fmt.Fprintf(sw, "    ;;\n")

	}
	fmt.Fprintf(w, script, sw.String())
	fmt.Fprintf(w, "compdef _goflag_complete %s\n", command)
}

// GenerateBashCompletions generates a Bash completion script for goflag.
func (ctx *Context) GenerateBashCompletions(w io.Writer) {
	fmt.Fprintln(w, functions)

	command := filepath.Base(os.Args[0])
	fmt.Fprintf(w, "\n# Bash completion script for %s\n", command)
	fmt.Fprintln(w, "_goflag_complete() {")
	fmt.Fprintln(w, "  COMPREPLY=()")

	// store variable to main program
	fmt.Fprintln(w, "  local program_name=\"${COMP_WORDS[0]}\"")
	fmt.Fprintln(w, "  local current_word=\"${COMP_WORDS[COMP_CWORD]}\"")
	fmt.Fprintln(w, "  local prev_word=\"${COMP_WORDS[COMP_CWORD - 1]}\"")

	//  Check if a subcommand already exists in words
	fmt.Fprintln(w, `
    subcommand_exists=0
    current_subcmd=""
    for word in "${COMP_WORDS[@]}"; do
        if [[  -n $word && "$word" != -* && "$word" != "--"* && "$word" != "$program_name" && "$word" != "$current_word" ]]; then
            # Found a non-flag word, assume it's a subcommand
            subcommand_exists=1
            current_subcmd="$word"
            break
        fi
    done
    `)

	// Generate completion options for global flags
	fmt.Fprint(w, "  local opts=(")

	for index, flag := range ctx.flags {
		if index > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, `"--%s"`, flag.name)
	}
	fmt.Fprint(w, ")\n")

	// declare an array of valid subcommands
	fmt.Fprint(w, "  local valid_subcommands=(")
	for index, cmd := range ctx.subcommands {
		if index > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, `"%s"`, cmd.name)
	}
	fmt.Fprintln(w, ")")

	for _, cmd := range ctx.subcommands {
		arrName := fmt.Sprintf("flags_%s", cmd.name)
		fmt.Fprintf(w, "  declare -a %s=(", arrName)

		// Generate completion options for subcommand flags
		for index, flag := range cmd.flags {
			if index > 0 {
				fmt.Fprint(w, " ")
			}
			fmt.Fprintf(w, "\"--%s\"", flag.name)
		}
		fmt.Fprintf(w, ")\n")
	}

	// Filter options based on the current subcommand
	fmt.Fprintln(w, `
    if [[ $subcommand_exists -eq 1 ]]; then  
        case "$current_subcmd" in
    `)

	for _, cmd := range ctx.subcommands {
		arrName := fmt.Sprintf("flags_%s", cmd.name)
		fmt.Fprintf(w, "          \"%s\")\n", cmd.name)
		fmt.Fprintf(w, "              COMPREPLY=( $(compgen -W \"${%s[*]}\" -- \"$current_word\") )\n", arrName)
		fmt.Fprintf(w, "          ;;\n")
	}

	fmt.Fprintln(w, `
            *)
                COMPREPLY=( $(compgen -W "${opts[*]} ${valid_subcommands[*]}" -- "$current_word") )
                ;;
        esac

	elif [[ "$prev_word" == "$program_name" ]]; then
        COMPREPLY=($(compgen -W "${opts[*]} ${valid_subcommands[*]}" -- "$current_word"))
    fi
    `)

	fmt.Fprintln(w, "}")

	fmt.Fprintln(w, "complete -F _goflag_complete "+command)
}

func GenerateCompletions(ctx Getter, cmd Getter) {
	zsh := cmd.GetBool("zsh")
	bash := cmd.GetBool("bash")
	outfile := cmd.GetString("out")

	if (zsh && bash) || (!zsh && !bash) {
		log.Fatalln("Specify one of zsh or bash")
	}

	var writer io.Writer = os.Stdout
	if outfile != "" {
		f, err := os.Create(outfile)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		writer = f
	}

	if zsh {
		ctx.(*Context).GenerateZshCompletions(writer)
	} else if bash {
		ctx.(*Context).GenerateBashCompletions(writer)
	}

}
