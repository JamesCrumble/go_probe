package probe_psql

const PSQL_QUERY string = `
with fraud_contractor_source as (
	select
		fraud_contractor_name as contractor_name,
		fraud_contractor_inn as contractor_inn,
		fc_extended.fraud_contractor_extended_status as contractor_status
	from fraud.fraud_contractor fc
	inner join fraud.fraud_contractor_extended fc_extended on fc_extended.fraud_contractor_id = fc.fraud_contractor_id
)
, contractors_to_delete as (
	select distinct contractor_name, contractor_inn
	from fraud_contractor_source
	where contractor_status not in ('unconfirmed', 'canceled')
)

select distinct
	contractor_name as contractor_name,
	array_remove(array_remove(contractor_inn, null), '') as contractor_inns
from fraud_contractor_source fc
`

type ContractorInfo struct {
	ContractorName string
	ContractorInns []string
}

func Probe() ([]ContractorInfo, error) {
	return probePsqlExtractDataDueReflect[ContractorInfo](PSQL_QUERY)
}
