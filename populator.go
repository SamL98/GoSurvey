package main

import (
	"log"
	"math/rand"
)

var currWave int
var res Response
var currRes Response

func GetResponse() {
	res = Response{wave: currWave}

	if err := pgManager.GetRandomResponse(&res); err != nil {
		log.Println("Error querying random row from postgres ", err)
	}

	/*if err := pgManager.MarkResponseAsUsed(res.id); err != nil {
		log.Println("Error marking id as used from postgres ", res.id, err)
	}*/

	PopulateQuestions()
}

func Shuffle(slc []Question) {
	for i := 1; i < len(slc); i++ {
		r := rand.Intn(i + 1)
		if i != r {
			slc[r], slc[i] = slc[i], slc[r]
		}
	}
}

func PopulateQuestions() {
	texts := [4]string{
		"According to a recent poll, <span id=\"s1\"></span>% of Americans say that they would prefer working under a male boss. <span id=\"s2\"></span>% of Americans would prefer to work under a female boss.",
		"Same sex marriage is a contested topic among Americans. In a poll conducted by the Pew Research Center, <span id=\"s1\"></span>% of respondents reported favoring same-sex marriage. <span id=\"s2\"></span>% reported opposing same-sex marriage.",
		"In 2007, <span id=\"s1\"></span> million Mexican immigrants lived in the United States. In 2014, <span id=\"s2\"></span> million Mexican immigrants lived in the US. Mexican immigrants have been at the center of one of the largest mass migrations in modern history.",
		"According to a database compiled by The Washington Post, 963 individuals where killed by police in 2016. Of those shot, <span id=\"s1\"></span> individuals were white and <span id=\"s2\"></span> were black.",
	}

	for i := range res.targets {
		res.targets[i].text = texts[i]
	}

	distractorTexts := [16]string{
		"Recently, the Pew Research Center ran a report about smartphone use in the United States. 64% of American adults were smartphone owners as of Spring 2015. This is a staggering number, when you realize that in 2011 that statistic was only 35%.",
		"The number of homeless in the United States stands at 564,708 people as of 2016. This study defined homelessness as: \"living on the streets, in cars, homeless shelters, or in subsidized transitional housing\". Of that over half a million people, about a quarter of them are children.",
		"According to PEW Research Center, gun ownership by households in the United States stands at 44%. This represents a 7-percentage point increase in the past two years. This is still lower than it was in the 1970s.",
		"According to preliminary reporting, the 2016 presidential election saw the lowest rate of voter turnout in 20 years. About 55.4% of eligible voters casted their ballots in the past election, which was down from about 60% in 2012",
		"The national debt currently is at $19.9 trillion. Although many believe a massive portion of the U.S. debt is owned by China, the country is second at $1.049 trillion to Japan, who owns $1.108 trillion American debt.",
		"A Gallup survey from 2016 revealed which countries Americans view most favorably. Canada tops the list. The next four favorites were Great Britain, France, Germany, and Japan. The bottom of the list includes Iran, Syria, and North Korea.",
		"A scene captured by a Google Street View vehicle in Kweneng, Botswana, that went viral seems to show the car hitting and possibly killing a donkey on the road. But Google insists the animal is fine.",
		"A Cheesecake Factory pasta dish with more than 3,000 calories - or more than a day and a half of the recommended caloric intake for an average adult - is one of the most unhealthy dishes at U.S. chain restaurants.",
		"An Indian company which plans to build a power plant that runs on city waste will aim for a $200 million IPO next month, according to a recent report from Bloomberg.",
		"If you are thinking about tweeting about clouds, pork, or exercise, think again. The Department of Homeland Security has been forced to release a list of keywords and phrases it uses to monitor various social networking sites.",
		"Vulture reports this morning that the Gremlins horror movie about creatures who can't be fed after midnight is perhaps on its way to a reboot. Steven Spielberg, who executive produced the original film has reportedly kept previous remake attempts from becoming reality.",
		"In recent months, the Food and Drug Administration has begun examining the safety of energy drinks following reports of several deaths and numerous injuries potentially associated with the products.",
		"The discovery of horse DNA in hamburgers on sale at supermarkets in Ireland and Britain is testing the appetite of meat lovers there. The Food Safety Authority of Ireland said that 10 out of 27 hamburger products it analyzed were found to contain horse DNA.",
		"Last year the US was ranked 10th place in the list of the World's Happiest Countries. This year the U.S.A. has slipped to 12th. This marks the first time that America is not in the top 10.",
		"A new study published in Administrative Science Quarterly shows how employees' wages change immediately after a male chief executive officer has a child. It found that when a male CEO has a child, his employees' wages decrease.",
		"Three major scientific projects set out this season to seek evidence of life in lakes deep under the Antarctic ice - evidence that could provide clues in the search for evidence of life elsewhere in the solar system such as under the surface of Enceladus, one of Saturn's moons.",
	}

	for i := range distractorTexts {
		res.distractors = append(res.distractors, Question{
			text:       distractorTexts[i],
			distractor: true,
		})
	}

	tmp := append(res.targets, res.distractors...)
	Shuffle(tmp)
	for i := range tmp {
		res.questions = append(res.questions, tmp[i])
	}
}
