package main

import (
	"image"

	"github.com/jim/kindleland"
)

var text = `In a village of La Mancha, the name of which I have no desire to call to mind, there lived not long since one of those gentlemen that keep a lance in the lance-rack, an old buckler, a lean hack, and a greyhound for coursing. An olla of rather more beef than mutton, a salad on most nights, scraps on Saturdays, lentils on Fridays, and a pigeon or so extra on Sundays, made away with three-quarters of his income. The rest of it went in a doublet of fine cloth and velvet breeches and shoes to match for holidays, while on week-days he made a brave figure in his best homespun. He had in his house a housekeeper past forty, a niece under twenty, and a lad for the field and market-place, who used to saddle the hack as well as handle the bill-hook. The age of this gentleman of ours was bordering on fifty; he was of a hardy habit, spare, gaunt-featured, a very early riser and a great sportsman. They will have it his surname was Quixada or Quesada (for here there is some difference of opinion among the authors who write on the subject), although from reasonable conjectures it seems plain that he was called Quexana. This, however, is of but little importance to our tale; it will be enough not to stray a hair’s breadth from the truth in the telling of it.

You must know, then, that the above-named gentleman whenever he was at leisure (which was mostly all the year round) gave himself up to reading books of chivalry with such ardour and avidity that he almost entirely neglected the pursuit of his field-sports, and even the management of his property; and to such a pitch did his eagerness and infatuation go that he sold many an acre of tillageland to buy books of chivalry to read, and brought home as many of them as he could get. But of all there were none he liked so well as those of the famous Feliciano de Silva’s composition, for their lucidity of style and complicated conceits were as pearls in his sight, particularly when in his reading he came upon courtships and cartels, where he often found passages like “the reason of the unreason with which my reason is afflicted so weakens my reason that with reason I murmur at your beauty;” or again, “the high heavens, that of your divinity divinely fortify you with the stars, render you deserving of the desert your greatness deserves.” Over conceits of this sort the poor gentleman lost his wits, and used to lie awake striving to understand them and worm the meaning out of them; what Aristotle himself could not have made out or extracted had he come to life again for that special purpose. He was not at all easy about the wounds which Don Belianis gave and took, because it seemed to him that, great as were the surgeons who had cured him, he must have had his face and body covered all over with seams and scars. He commended, however, the author’s way of ending his book with the promise of that interminable adventure, and many a time was he tempted to take up his pen and finish it properly as is there proposed, which no doubt he would have done, and made a successful piece of work of it too, had not greater and more absorbing thoughts prevented him.

Many an argument did he have with the curate of his village (a learned man, and a graduate of Siguenza) as to which had been the better knight, Palmerin of England or Amadis of Gaul. Master Nicholas, the village barber, however, used to say that neither of them came up to the Knight of Phoebus, and that if there was any that could compare with him it was Don Galaor, the brother of Amadis of Gaul, because he had a spirit that was equal to every occasion, and was no finikin knight, nor lachrymose like his brother, while in the matter of valour he was not a whit behind him. In short, he became so absorbed in his books that he spent his nights from sunset to sunrise, and his days from dawn to dark, poring over them; and what with little sleep and much reading his brains got so dry that he lost his wits. His fancy grew full of what he used to read about in his books, enchantments, quarrels, battles, challenges, wounds, wooings, loves, agonies, and all sorts of impossible nonsense; and it so possessed his mind that the whole fabric of invention and fancy he read of was true, that to him no history in the world had more reality in it. He used to say the Cid Ruy Diaz was a very good knight, but that he was not to be compared with the Knight of the Burning Sword who with one back-stroke cut in half two fierce and monstrous giants. He thought more of Bernardo del Carpio because at Roncesvalles he slew Roland in spite of enchantments, availing himself of the artifice of Hercules when he strangled Antaeus the son of Terra in his arms. He approved highly of the giant Morgante, because, although of the giant breed which is always arrogant and ill-conditioned, he alone was affable and well-bred. But above all he admired Reinaldos of Montalban, especially when he saw him sallying forth from his castle and robbing everyone he met, and when beyond the seas he stole that image of Mahomet which, as his history says, was entirely of gold. To have a bout of kicking at that traitor of a Ganelon he would have given his housekeeper, and his niece into the bargain.`

func main() {
	tv := kindleland.NewTextView(text, image.Rect(50, 50, 550, 750))

	fb, err := kindleland.NewFrameBuffer("/dev/fb0", 600, 800)
	if err != nil {
		panic(err)
	}

	fb.ApplyImage(tv.Render())

	err = fb.UpdateScreen()
	if err != nil {
		panic(err)
	}
}
