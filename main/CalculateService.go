package main

type CalculateService interface {
	CalculateLength(fact HumanFact) HumanFact
}

func NewCalculateService() CalculateService {
	return HumanFact{}
}

func (h HumanFact) CalculateLength(fact HumanFact) HumanFact {
	return HumanFact{
		fact.Fact,
		len(fact.Fact) / 2,
	}
}
