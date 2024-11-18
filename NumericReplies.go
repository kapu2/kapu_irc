package main

const (
	RPL_WELCOME = string("001") // done

	RPL_LUSERCLIENT   = string("251") // done
	RPL_LUSEROP       = string("252") // done
	RPL_LUSERUNKNOWN  = string("253") // done
	RPL_LUSERCHANNELS = string("254") // done
	RPL_LUSERME       = string("255") // done
	RPL_ADMINME       = string("256") // done
	RPL_ADMINLOC1     = string("257") // done
	RPL_ADMINLOC2     = string("258") // done
	RPL_ADMINEMAIL    = string("259") // done
	RPL_TRYAGAIN      = string("263") // done
	RPL_LOCALUSERS    = string("265") // done
	RPL_GLOBALUSERS   = string("266") // done

	RPL_TOPIC        = string("332") // done
	RPL_TOPICWHOTIME = string("333") // done
	RPL_NAMREPLY     = string("353") // done
	RPL_ENDOFNAMES   = string("366") // done
	RPL_MOTD         = string("372") // done
	RPL_MOTDSTART    = string("375") // done
	RPL_ENDOFMOTD    = string("376") // done
)
