package picasa;

// globals
const (
	PICASA_INI_FILE		= ".picasa.ini"
	NOMEDIA_FILE		= ".nomedia"
	PICASA_VDIR_prefix	= "!"
	PICASA_VDIR_STARRED	= "STARRED"
)

// other presets
const (
	CACHE_EXPIRE		= 1800
)


// rules ident
const (
	rule_MATCH_FILE			= "(.*\\..*)"

	RULE_IDENT_ROOT			= 1
	RULE_MATCH_ROOT			= "^$"

	RULE_IDENT_ROOT_STARRED_DIR	= 10
	RULE_MATCH_ROOT_STARRED_DIR	= "^" + PICASA_VDIR_prefix + PICASA_VDIR_STARRED + "$"

	RULE_IDENT_ROOT_STARRED_YEAR	= 11
	RULE_MATCH_ROOT_STARRED_YEAR	= "^" + PICASA_VDIR_prefix + PICASA_VDIR_STARRED + "/(\\d{4})$"

	RULE_IDENT_ROOT_STARRED_YEAR_FILE	= 12
	RULE_MATCH_ROOT_STARRED_YEAR_FILE	= "^" + PICASA_VDIR_prefix + PICASA_VDIR_STARRED + "/\\d{4}/" + rule_MATCH_FILE + "$"

	RULE_IDENT_SUBDIR_STARRED_DIR	= 100
	RULE_MATCH_SUBDIR_STARRED_DIR	= "^(.*)/" + PICASA_VDIR_prefix + PICASA_VDIR_STARRED + "$"

	RULE_IDENT_SUBDIR_STARRED_FILE	= 101
	RULE_MATCH_SUBDIR_STARRED_FILE	= "^(.*)/" + PICASA_VDIR_prefix + PICASA_VDIR_STARRED + "/" + rule_MATCH_FILE + "$"
)