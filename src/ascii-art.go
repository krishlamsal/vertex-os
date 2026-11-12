package main

import "strings"

var appArt = map[string]string{
	"Calculator": `
###############
#             #
# ########### #
# #         # #
# ########### #
#             #
# ### ### ### #
# ### ### ### #
#             #
# ### ### ### #
# ### ### # # #
#         # # #
# ### ### # # #
# ### ### ### #
#             #
###############
`,
	"Notes": `
 ############
 #          #
 #          #
 # # ###### #
 #          #
 # # ###### #
 #          #
 # # ###### #
 #          #
 # # ###### #
 #          #
 # # ###### #
 #          #
 #          #
 ############
`,
	"Calendar": `

  #         #
###############
#             #
#  #  #  #  # #
#             #
###############
#             #
#       ##### #
#       ## ## #
#       ## ## #
#       ##### #
#             #
###############
`,
	"Camera": `
       ####
      #    #
  ##############
 ################
 #####      #####
 ####  ####  ####
 ####  ####  ####
 #####      #####
 ################
 #              #
 #              #
 ################
`,
	"Clock": `
      ####
    ##    ##
   #        #
  #    ##    #
 #     ##     #
 #     ##     #
 #     ###### #
 #            #
  #          #
   #        #
    ##    ##
      ####
`,
	"Stopwatch": `
    ########
      ####
    ##    ##
   #   ##   #
  #    ##    #
 #    ####    #
 #    ####    #
 #            #
  #          #
   #        #
    ##    ##
      ####
`,
	"Timer": `
   ##########
   #        #
   #  ####  #
    #  ##  #
     ##  ##
       ##
       ##
     ##  ##
    #      #
   #        #
   #        #
   ##########
`,
	"Weather": `
       ###
      #   #
  ####     #
 #    # ###
#      #   #
#          #
#          #
############
`,
	"Browser": `
      #####
    ##  #  ##
   # #  #  # #
  #############
 #  #   #   #  #
 #  #   #   #  #
 ###############
 #  #   #   #  #
  ####  #  ####
   # ####### #
    ##  #  ##
      #####
`,
	"Phone": `
##########
#        #
# ###### #
# #    # #
# #    # #
# #    # #
# #    # #
# ###### #
#        #
#   ##   #
#        #
##########
`,
	"Phonebook": `
   ##########
  #          #
 #    ####    #
 #   #    #   #
 #   #    #   #
 #    ####    #
 #   ######   #
 #  ########  #
 #  ########  #
 #            #
  #          #
   ##########
`,
	"Messages": `
  ############
 #            #
 #     ####   #
 #            #
 #   ####     #
 #            #
  ############
         ##
        ##
       #
`,
	"Files": `
  #######
 ##############
 #            #
 #            #
 #            #
 #            #
 #            #
 #            #
 #            #
 #            #
 #            #
 #            #
 ##############
`,
	"Text Editor": `
  ############
  #          #
  # ######## #
  # #  ##  # #
  #    ##    #
  #    ##    #
  #    ##    #
  #    ##    #
  #    ##    #
  #   ####   #
  #          #
  ############
`,
	"Image Viewer": `
 ##############
 #            #
 #   ##       #
 #   ##       #
 #         #  #
 #        ### #
 #   #   #### #
 #  ### ##### #
 # ########## #
 #            #
 #            #
 ##############
`,
	"Audio Player": `
    ##########
    #        #
    #        #
    #        #
    #        #
    #        #
    #        #
 ####     ####
#####    #####
#####    #####
 ###      ###

`,
	"Settings": `
##############
#            #
#            #
# ###        #
# ########## #
# ###        #
#        ### #
# ########## #
#        ### #
#            #
#            #
##############
`,
	"Help": `
  ######
 ##    ##
##      ##
        ##
       ##
      ##
    ##
    ##
    ##

    ##
    ##
`,
	"System Monitor": `
 ##############
 ##           #
 ###          #
 # ##     #   #
 # ##    ###  #
 #  ##  ## ## #
 #  ## ##  ## #
 #   ###    ###
 #            #
 ##############
     ######
     ######
`,
	"Terminal": `
################
#              #
#              #
#              #
#              #
#  #           #
#   #          #
#    #         #
#   #          #
#  #  ######## #
#              #
################
`,
}

var numberArt = map[string]string{
	"1": `
##
 #
 #
 #
 #
###
`,
	"2": `
###
  #
###
#
#
###
`,
	"3": `
###
  #
###
  #
  #
###
`,
	"4": `
# #
# #
# #
###
  #
  #
`,
	"5": `
###
#
###
  #
  #
###
`,
	"6": `
###
#
###
# #
# #
###
`,
	"7": `
###
# #
  #
  #
  #
  #
`,
	"8": `
###
# #
###
# #
# #
###
`,
	"9": `
###
# #
# #
###
  #
###
`,
	"0": `
###
# #
# #
# #
# #
###
`,
	":": `
#
#


#
#
`, ".": `
.
.
.
.
#
#
`,
}

// Converts a sprite into half-block cell rows, padded and centered. Returns nil if no sprite.
func artRows(name string, list map[string]string, maxWidth int) []string {
	raw, ok := list[name]
	if !ok {
		return nil
	}

	lines := strings.Split(raw, "\n")

	// Trim leading/trailing blank lines from the raw strings
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return nil
	}

	width := 0
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
		if w := len(lines[i]); w > width {
			width = w
		}
	}

	// Center crop sprites that are too wide
	if width > maxWidth {
		start := (width - maxWidth) / 2
		for i, line := range lines {
			switch {
			case len(line) > start+maxWidth:
				lines[i] = line[start : start+maxWidth]
			case len(line) > start:
				lines[i] = line[start:]
			default:
				lines[i] = ""
			}
		}
		width = maxWidth
	}

	// Right pad lines to the sprite
	for i, line := range lines {
		lines[i] = line + strings.Repeat(" ", width-len(line))
	}

	// Center horizontally
	pad := maxWidth - width
	left := pad / 2
	right := pad - left
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
	}

	var rows []string
	for y := 0; y < len(lines); y += 2 {
		top := lines[y]
		bottom := ""
		if y+1 < len(lines) {
			bottom = lines[y+1]
		}

		var b strings.Builder
		for x := 0; x < maxWidth; x++ {
			topOn := top[x] == '#'
			bottomOn := x < len(bottom) && bottom[x] == '#'
			switch {
			case topOn && bottomOn:
				b.WriteRune('█')
			case topOn:
				b.WriteRune('▀')
			case bottomOn:
				b.WriteRune('▄')
			default:
				b.WriteRune(' ')
			}
		}
		rows = append(rows, b.String())
	}

	return rows
}
