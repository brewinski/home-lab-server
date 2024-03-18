package cfg

const (
	VQBaseUrl = "https://vqmetro24s1.softr.app/v1/integrations/airtable/dc83c433-262d-48a0-915f-2cf124cceeb8/app4eDFcW0KK8A7xt"
	// VQ Ladder specific variables (change between seasons)
	VQLadderPath            = "/Ladder/records?block_id=4cf2b9cc-8241-4332-9df5-47a68e375c5a"
	VQLadderPageID          = "a7233511-bb1c-4840-9f02-d3198caf05f4"
	VQLadderFilterByFormula = "(LOWER(\"MD\") = LOWER(ARRAYJOIN({Division})))"
)
