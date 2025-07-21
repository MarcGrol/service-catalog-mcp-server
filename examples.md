# Examples request-response pairs

Start the in stdio mode

    learnmcp

And type json rpc in the command-line

## suggest_candidates

**request:**

    {"method":"tools/call","params":{"name":"suggest_candidates","arguments":{"keyword":"partner"}},"jsonrpc":"2.0","id":9}

**response:**

    {"jsonrpc":"2.0","id":9,"result":{"content":[{"type":"text","text":"{\n  \"status\": \"success\",\n  \"data\": {\n    \"Modules\": [\n      \"partner\",\n      \"partner-jobs\",\n      \"common/partner\",\n      \"ui/resources/partner\",\n      \"communication/services/partner\"\n    ],\n    \"Teams\": [\n      \"PartnerExperience\",\n      \"PartnerExperience_FE\",\n      \"PlatformIntegrationExperience\"\n    ],\n    \"Interfaces\": [\n      \"PartnerTermsResourceV1\",\n      \"PartnerReferralResourceV1\",\n      \"PartnerMarketingResourceV1\",\n      \"PartnerDocumentsResourceV1\",\n      \"PartnerOnboardingResourceV1\"\n    ],\n    \"Databases\": [\n      \"partner\"\n    ]\n  }\n}"}]}}

## get_interface

**request:**

    {"method":"tools/call","params":{"name":"get_interface","arguments":{"interface_id":"com.adyen.services.partner.PartnerCommissionService"}},"jsonrpc":"2.0","id":32}

**response:**

    {"jsonrpc":"2.0","id":32,"result":{"content":[{"type":"text","text":"{\n  \"status\": \"success\",\n  \"data\": {\n    \"InterfaceID\": \"com.adyen.services.partner.PartnerCommissionService\",\n    \"Description\": \"PartnerCommissionService\",\n    \"Kind\": \"RPL\",\n    \"Spec\": \"partner/src/main/resources/rpl-partner-commission.xml\",\n    \"MethodCount\": 14,\n    \"Methods\": [\n      \"addPartnerCommissionMerchants\",\n      \"approveCommissionStatement\",\n      \"deleteAllPartnerCommissionMerchants\",\n      \"deletePartnerCommissionMerchant\",\n      \"getPartnerStatementDetails\",\n      \"getPartnerStatements\",\n      \"getPartnerStatementsUnapproved\",\n      \"listAccountGroupMembers\",\n      \"listCompaniesForConnection\",\n      \"listMerchantAccountGroups\",\n      \"listPartnerCommissionCompanies\",\n      \"listPartnerCommissionMerchants\",\n      \"listUngroupedMerchants\",\n      \"recalculateStatements\"\n    ]\n  }\n}"}]}}

