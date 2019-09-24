'''
The job of this script is to convert the output of the Go job into 
what is required to fit whatever chart libary is being used and to 
possibly merge related sets together
'''

import json
import os
import operator


def filesPerProject():
    '''
    Converts the output of filesPerProject into something
    we can throw into a chart library since it needs to 
    be sorted
    It is a count of the number of projects that have a number of files

    EG. files:project where 123 projects have 2 files in them
    https://jsfiddle.net/uLw08scq/
    '''
    data = '[]'
    with open('./results/filesPerProject.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([int(x), y])

    def cmp(a, b):
        if a == b:
            return 0
        if a < b:
            return -1
        return 1

    new.sort(cmp)

    with open("./results/filesPerProject_converted.json", "w") as text_file:
        text_file.write(json.dumps(new, sort_keys=True))


def filesPerProjectPercentile():
    '''
    Want to plot by percentile
    https://jsfiddle.net/61mv2xrc/
    https://jsfiddle.net/61mv2xrc/1/


    google.charts.load('current', {packages: ['corechart', 'line']});
google.charts.setOnLoadCallback(drawBasic);

function drawBasic() {

      var data = new google.visualization.DataTable();
      data.addColumn('number', 'Day');
      data.addColumn('number', 'Files');
      /* data.addColumn('number', 'The Avengers') */;
      /* data.addColumn('number', 'Transformers: Age of Extinction') */;

      data.addRows(
      [[1, 6.206737425011537], [2, 12.7365020766036], [3, 16.958929395477618], [4, 22.81956622058145], [5, 26.09598523304107], [6, 28.72634979233964], [7, 31.541301338255654], [8, 33.82556529764652], [9, 36.31748961698201], [10, 38.44023996308261], [11, 40.216889709275506], [12, 42.13197969543148], [13, 43.70096908167975], [14, 45.59298569450855], [15, 47.485002307337346], [16, 48.80018458698663], [17, 49.86155976003693], [18, 50.4845408398708], [19, 51.75357637286573], [20, 52.51499769266268], [21, 53.62251961236734], [22, 54.47623442547302], [23, 55.16843562528843], [24, 56.13751730503001], [25, 56.8527918781726], [26, 57.36040609137057], [27, 57.84494693124136], [28, 58.49100138440241], [29, 58.99861559760038], [30, 59.64467005076143], [31, 60.22150438394094], [32, 60.70604522381173], [33, 61.513613290263045], [34, 62.02122750346101], [35, 62.55191508998616], [36, 62.944162436548226], [37, 63.428703276419014], [38, 63.84402399630826], [39, 64.16705122288879], [40, 64.46700507614213], [41, 64.9746192893401], [42, 65.38994000922935], [43, 65.75911398246424], [44, 66.01292108906323], [45, 66.26672819566221], [46, 66.54360867558837], [47, 67.00507614213198], [48, 67.3281033687125], [49, 67.78957083525611], [50, 68.02030456852792], [51, 69.98154130133825], [52, 71.15828334102446], [53, 73.05029995385325], [54, 73.55791416705122], [55, 73.92708814028612], [56, 74.15782187355792], [57, 74.38855560682973], [58, 74.82694970004616], [59, 75.10383017997232], [60, 75.38071065989848], [61, 75.49607752653438], [62, 75.8191047531149], [63, 76.1652053530226], [64, 76.34979233964005], [65, 76.46515920627594], [66, 76.60359944623903], [67, 76.81125980618366], [68, 76.9958467928011], [69, 77.08814028610982], [70, 77.24965389940009], [71, 77.43424088601753], [72, 77.54960775265343], [73, 77.78034148592523], [74, 77.89570835256113], [75, 78.03414859252422], [76, 78.26488232579602], [77, 78.44946931241347], [78, 78.61098292570374], [79, 78.72634979233963], [80, 78.93401015228426], [81, 79.02630364559298], [82, 79.16474388555606], [83, 79.18781725888324], [84, 79.32625749884633], [85, 79.51084448546378], [86, 79.78772496538994], [87, 79.97231195200739], [88, 80.04153207198893], [89, 80.22611905860637], [90, 80.34148592524227], [91, 80.52607291185971], [92, 80.64143977849561], [93, 80.73373327180433], [94, 80.82602676511304], [95, 80.91832025842176], [96, 80.96446700507612], [97, 81.07983387171201], [98, 81.17212736502073], [99, 81.28749423165662], [100, 81.37978772496534], [101, 81.54130133825561], [102, 81.6566682048915], [103, 81.7720350715274], [104, 81.79510844485458], [105, 81.8874019381633], [106, 82.04891555145356], [107, 82.21042916474383], [108, 82.32579603137972], [109, 82.57960313797871], [110, 82.64882325796025], [111, 82.76419012459614], [112, 82.87955699123204], [113, 82.90263036455922], [114, 82.9257037378864], [115, 82.94877711121357], [116, 83.01799723119511], [117, 83.1564374711582], [118, 83.20258421781256], [119, 83.29487771112127], [120, 83.41024457775717], [121, 83.45639132441153], [122, 83.52561144439306], [123, 83.57175819104742], [124, 83.64097831102896], [125, 83.75634517766485], [126, 83.80249192431921], [127, 83.89478541762793], [128, 83.96400553760947], [129, 84.01015228426382], [130, 84.14859252422691], [131, 84.19473927088127], [132, 84.26395939086281], [133, 84.31010613751717], [134, 84.42547300415306], [136, 84.51776649746178], [137, 84.54083987078896], [138, 84.63313336409767], [139, 84.70235348407921], [140, 84.74850023073357], [141, 84.8177203507151], [142, 84.933087217351], [143, 85.00230733733254], [144, 85.02538071065972], [146, 85.11767420396843], [147, 85.18689432394997], [148, 85.20996769727715], [149, 85.23304107060433], [150, 85.27918781725869], [151, 85.32533456391305], [152, 85.34840793724022], [153, 85.39455468389458], [154, 85.64836179049357], [155, 85.76372865712946], [156, 85.85602215043818], [157, 85.90216889709254], [158, 85.99446239040125], [159, 86.08675588370997], [160, 86.17904937701869], [162, 86.2713428703274], [163, 86.31748961698176], [165, 86.36363636363612], [166, 86.3867097369633], [167, 86.45592985694483], [168, 86.47900323027201], [169, 86.50207660359919], [174, 86.52514997692637], [175, 86.54822335025355], [176, 86.57129672358073], [177, 86.61744347023509], [178, 86.75588371019818], [179, 86.82510383017971], [180, 86.84817720350689], [181, 86.89432395016125], [183, 86.91739732348843], [184, 86.96354407014279], [186, 86.98661744346997], [187, 87.0558375634515], [188, 87.07891093677868], [189, 87.12505768343304], [191, 87.1712044300874], [192, 87.19427780341458], [193, 87.21735117674176], [194, 87.26349792339612], [195, 87.2865712967233], [197, 87.30964467005047], [198, 87.37886479003201], [200, 87.42501153668637], [201, 87.44808491001355], [203, 87.47115828334073], [206, 87.4942316566679], [207, 87.51730502999509], [208, 87.58652514997662], [210, 87.63267189663098], [212, 87.67881864328534], [213, 87.70189201661252], [214, 87.7249653899397], [217, 87.74803876326688], [218, 87.79418550992123], [219, 87.84033225657559], [221, 87.86340562990277], [222, 87.90955237655713], [223, 87.95569912321149], [227, 88.00184586986585], [231, 88.0479926165202], [232, 88.09413936317456], [233, 88.11721273650174], [234, 88.1633594831561], [237, 88.18643285648328], [238, 88.23257960313764], [239, 88.25565297646482], [241, 88.32487309644635], [242, 88.34794646977353], [244, 88.37101984310071], [245, 88.44023996308225], [246, 88.46331333640943], [247, 88.50946008306379], [248, 88.55560682971814], [249, 88.57868020304532], [250, 88.64790032302686], [251, 88.69404706968122], [252, 88.7171204430084], [254, 88.76326718966276], [255, 88.78634056298993], [256, 88.80941393631711], [258, 88.85556068297147], [259, 88.87863405629865], [260, 88.90170742962583], [263, 88.92478080295301], [265, 88.94785417628019], [267, 88.97092754960737], [268, 88.99400092293455], [272, 89.01707429626173], [277, 89.0401476695889], [278, 89.06322104291608], [280, 89.08629441624326], [281, 89.10936778957044], [284, 89.1555145362248], [285, 89.22473465620634], [286, 89.24780802953352], [287, 89.2708814028607], [288, 89.29395477618787], [289, 89.40932164282377], [290, 89.43239501615095], [292, 89.45546838947813], [293, 89.4785417628053], [294, 89.50161513613249], [296, 89.54776188278684], [297, 89.5939086294412], [298, 89.61698200276838], [299, 89.64005537609556], [300, 89.7092754960771], [302, 89.73234886940428], [303, 89.75542224273146], [310, 89.77849561605863], [315, 89.96308260267608], [321, 89.98615597600326], [322, 90.00922934933044], [324, 90.07844946931198], [326, 90.10152284263916], [329, 90.12459621596633], [331, 90.14766958929351], [335, 90.1707429626207], [336, 90.19381633594787], [337, 90.21688970927505], [341, 90.23996308260223], [342, 90.28610982925659], [346, 90.30918320258377], [347, 90.33225657591095], [349, 90.35532994923813], [353, 90.3784033225653], [354, 90.7475772958002], [356, 90.77065066912738], [357, 90.79372404245456], [359, 90.81679741578174], [360, 90.90909090909045], [361, 90.93216428241763], [365, 90.95523765574481], [366, 90.97831102907199], [370, 91.00138440239917], [371, 91.04753114905353], [372, 91.0706045223807], [373, 91.09367789570788], [377, 91.13982464236224], [378, 91.1859713890166], [381, 91.20904476234378], [383, 91.25519150899814], [385, 91.27826488232532], [390, 91.32441162897968], [395, 91.34748500230685], [397, 91.37055837563403], [398, 91.43977849561557], [399, 91.46285186894275], [401, 91.53207198892429], [402, 91.624365482233], [403, 91.67051222888736], [405, 91.69358560221454], [406, 91.78587909552326], [407, 91.83202584217761], [409, 91.8550992155048], [413, 91.87817258883197], [419, 91.90124596215915], [420, 91.94739270881351], [422, 91.97046608214069], [423, 91.99353945546787], [424, 92.01661282879505], [427, 92.0627595754494], [428, 92.08583294877658], [429, 92.10890632210376], [431, 92.13197969543094], [436, 92.15505306875812], [438, 92.1781264420853], [439, 92.22427318873966], [440, 92.24734656206684], [452, 92.2934933087212], [457, 92.31656668204838], [462, 92.33964005537555], [464, 92.36271342870273], [472, 92.38578680202991], [473, 92.40886017535709], [474, 92.43193354868427], [476, 92.45500692201145], [481, 92.47808029533863], [482, 92.50115366866581], [483, 92.52422704199299], [488, 92.54730041532017], [492, 92.57037378864734], [493, 92.59344716197452], [500, 92.6165205353017], [505, 92.63959390862888], [508, 92.66266728195606], [510, 92.68574065528324], [514, 92.7318874019376], [516, 92.75496077526478], [518, 92.77803414859196], [519, 92.80110752191914], [520, 92.82418089524631], [524, 92.8472542685735], [525, 92.87032764190067], [527, 92.89340101522785], [531, 92.91647438855503], [532, 92.93954776188221], [533, 92.96262113520939], [535, 92.98569450853657], [539, 93.00876788186375], [540, 93.03184125519093], [543, 93.0549146285181], [544, 93.07798800184528], [550, 93.10106137517246], [551, 93.12413474849964], [553, 93.14720812182682], [554, 93.170281495154], [559, 93.19335486848118], [564, 93.21642824180836], [567, 93.26257498846272], [582, 93.2856483617899], [591, 93.30872173511708], [592, 93.33179510844425], [593, 93.35486848177143], [599, 93.37794185509861], [608, 93.40101522842579], [612, 93.42408860175297], [618, 93.51638209506169], [619, 93.53945546838887], [620, 93.65482233502476], [621, 93.70096908167912], [622, 93.74711582833348], [623, 93.79326257498784], [624, 93.81633594831501], [627, 93.8394093216422], [632, 93.86248269496937], [635, 93.88555606829655], [639, 93.90862944162373], [648, 93.93170281495091], [649, 93.95477618827809], [650, 94.0470696815868], [652, 94.07014305491398], [657, 94.09321642824116], [668, 94.11628980156834], [672, 94.13936317489552], [673, 94.1624365482227], [679, 94.18550992154988], [686, 94.20858329487706], [687, 94.23165666820424], [689, 94.25473004153142], [695, 94.2778034148586], [696, 94.30087678818578], [700, 94.34702353484013], [706, 94.37009690816731], [707, 94.39317028149449], [710, 94.41624365482167], [712, 94.43931702814885], [719, 94.46239040147603], [726, 94.48546377480321], [728, 94.50853714813039], [730, 94.53161052145757], [731, 94.55468389478474], [732, 94.57775726811192], [736, 94.6008306414391], [743, 94.62390401476628], [745, 94.64697738809346], [750, 94.69312413474782], [756, 94.716197508075], [757, 94.73927088140218], [763, 94.76234425472936], [764, 94.78541762805654], [769, 94.80849100138371], [770, 94.8315643747109], [775, 94.85463774803807], [778, 94.87771112136525], [786, 94.92385786801961], [795, 94.94693124134679], [796, 94.97000461467397], [797, 94.99307798800115], [798, 95.01615136132833], [799, 95.0392247346555], [801, 95.06229810798268], [816, 95.08537148130986], [820, 95.10844485463704], [822, 95.13151822796422], [826, 95.1545916012914], [827, 95.17766497461858], [829, 95.20073834794576], [830, 95.22381172127294], [844, 95.24688509460012], [847, 95.2699584679273], [862, 95.29303184125448], [871, 95.33917858790883], [874, 95.36225196123601], [882, 95.38532533456319], [885, 95.40839870789037], [889, 95.45454545454473], [890, 95.47761882787191], [899, 95.50069220119909], [904, 95.52376557452627], [908, 95.54683894785344], [923, 95.56991232118062], [925, 95.5929856945078], [939, 95.61605906783498], [948, 95.63913244116216], [955, 95.66220581448934], [956, 95.68527918781652], [963, 95.7083525611437], [976, 95.73142593447088], [983, 95.75449930779806], [993, 95.77757268112524], [1004, 95.80064605445241], [1006, 95.8237194277796], [1008, 95.86986617443395], [1017, 95.89293954776113], [1020, 95.91601292108831], [1021, 95.93908629441549], [1033, 95.96215966774267], [1046, 95.98523304106985], [1058, 96.00830641439703], [1059, 96.0313797877242], [1061, 96.05445316105138], [1069, 96.07752653437856], [1073, 96.10059990770574], [1080, 96.12367328103292], [1111, 96.1467466543601], [1136, 96.16982002768728], [1145, 96.19289340101446], [1161, 96.21596677434164], [1165, 96.262113520996], [1170, 96.28518689432318], [1174, 96.30826026765035], [1178, 96.33133364097753], [1185, 96.35440701430471], [1225, 96.37748038763189], [1255, 96.40055376095907], [1256, 96.42362713428625], [1319, 96.44670050761343], [1441, 96.46977388094061], [1444, 96.51592062759497], [1447, 96.53899400092214], [1466, 96.56206737424932], [1470, 96.5851407475765], [1476, 96.60821412090368], [1479, 96.63128749423086], [1515, 96.65436086755804], [1524, 96.67743424088522], [1525, 96.7005076142124], [1554, 96.72358098753958], [1583, 96.74665436086676], [1640, 96.76972773419394], [1644, 96.79280110752111], [1658, 96.8158744808483], [1666, 96.83894785417547], [1674, 96.86202122750265], [1691, 96.88509460082983], [1712, 96.90816797415701], [1742, 96.93124134748419], [1765, 96.95431472081137], [1775, 96.97738809413855], [1780, 97.00046146746573], [1810, 97.0235348407929], [1872, 97.04660821412008], [1913, 97.06968158744726], [1919, 97.09275496077444], [1922, 97.11582833410162], [1938, 97.16197508075598], [2019, 97.18504845408316], [2070, 97.20812182741034], [2145, 97.23119520073752], [2153, 97.2542685740647], [2169, 97.27734194739188], [2190, 97.30041532071905], [2207, 97.32348869404623], [2217, 97.34656206737341], [2225, 97.36963544070059], [2247, 97.39270881402777], [2262, 97.41578218735495], [2263, 97.43885556068213], [2294, 97.46192893400931], [2311, 97.48500230733649], [2392, 97.50807568066367], [2411, 97.53114905399084], [2425, 97.55422242731802], [2440, 97.5772958006452], [2464, 97.60036917397238], [2536, 97.62344254729956], [2546, 97.64651592062674], [2603, 97.66958929395392], [2677, 97.6926626672811], [2686, 97.71573604060828], [2779, 97.73880941393546], [2796, 97.76188278726264], [2921, 97.78495616058981], [2929, 97.808029533917], [2995, 97.83110290724417], [3036, 97.85417628057135], [3050, 97.87724965389853], [3063, 97.90032302722571], [3067, 97.92339640055289], [3098, 97.94646977388007], [3099, 97.96954314720725], [3160, 97.99261652053443], [3181, 98.0156898938616], [3206, 98.03876326718878], [3221, 98.06183664051596], [3266, 98.08491001384314], [3351, 98.10798338717032], [3362, 98.1310567604975], [3386, 98.38486386709648], [3387, 98.45408398707802], [3389, 98.4771573604052], [3462, 98.50023073373238], [3483, 98.52330410705956], [3525, 98.54637748038674], [3567, 98.56945085371392], [3573, 98.5925242270411], [3593, 98.63867097369545], [3596, 98.66174434702263], [3598, 98.68481772034981], [3896, 98.70789109367699], [3962, 98.73096446700417], [4203, 98.75403784033135], [4267, 98.77711121365853], [4338, 98.80018458698571], [4836, 98.82325796031289], [4923, 98.84633133364007], [5175, 98.86940470696725], [5323, 98.89247808029442], [5359, 98.9155514536216], [5464, 98.93862482694878], [5769, 98.96169820027596], [5908, 98.98477157360314], [6497, 99.00784494693032], [6605, 99.0309183202575], [6998, 99.05399169358468], [7000, 99.07706506691186], [7148, 99.10013844023904], [7640, 99.12321181356621], [7941, 99.1462851868934], [8912, 99.16935856022057], [9349, 99.19243193354775], [9387, 99.21550530687493], [9864, 99.23857868020211], [10214, 99.26165205352929], [10330, 99.28472542685647], [10801, 99.30779880018365], [10882, 99.33087217351083], [10950, 99.353945546838], [11325, 99.37701892016518], [12192, 99.40009229349236], [13152, 99.42316566681954], [14011, 99.44623904014672], [14455, 99.4693124134739], [15899, 99.49238578680108], [20533, 99.51545916012826], [24614, 99.53853253345544], [31601, 99.56160590678262], [32464, 99.5846792801098], [32910, 99.60775265343698], [34354, 99.63082602676415], [34370, 99.65389940009133], [34377, 99.67697277341851], [34682, 99.70004614674569], [35170, 99.72311952007287], [35241, 99.74619289340005], [35531, 99.76926626672723], [35652, 99.79233964005441], [35838, 99.81541301338159], [36580, 99.83848638670877], [36674, 99.86155976003595], [36681, 99.88463313336312], [39255, 99.9077065066903], [47281, 99.93077988001748], [61645, 99.95385325334466], [87456, 99.97692662667184], [97247, 99.99999999999902]]
      );

      var options = {
        chart: {
          title: 'Box Office Earnings in First Two Weeks of Opening',
          subtitle: 'in millions of dollars (USD)'
        },
        width: 900,
        height: 500
      };

      
      
      var chart = new google.visualization.LineChart(document.getElementById('chart_div'));

      chart.draw(data, options);
    }
    '''
    data = '[]'
    with open('./results/filesPerProject.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([int(x), y])

    def cmp(a, b):
        if a == b:
            return 0
        if a < b:
            return -1
        return 1

    new.sort(cmp)

    # total number of projects
    projectCount = 0
    for x in new:
        projectCount += x[1]


    # Percentages
    percent = []
    for x in new:
        if len(percent) == 0:
            percent.append([x[0], (float(x[1])/float(projectCount)) * float(100)])
        else:
            percent.append([x[0], (float(x[1])/float(projectCount)) * float(100) + percent[len(percent)-1][1]])

    percent_log = []

    # count 1 to 10
    # count by 10 to 100
    # count by 100 to 1000
    # count by 1000 to 10000
    # count by 10000 to 100000

    tens = {}
    for x in percent:
        if x[0] <= 10:
            percent_log.append(x)

        if x[0] > 10 and x[0] <= 100:
            # rount to nearest value of 10
            r = int(roundup(x[0], 10.0))
            tens[r] = x[1]

        if x[0] > 100 and x[0] <= 1000:
            # rount to nearest value of 100
            r = int(roundup(x[0], 100.0))
            tens[r] = x[1]

        # if x[0] > 1000 and x[0] <= 10000:
        #     # rount to nearest value of 100
        #     r = int(roundup(x[0], 1000.0))
        #     tens[r] = x[1]

        # if x[0] > 10000 and x[0] <= 100000:
        #     # rount to nearest value of 100
        #     r = int(roundup(x[0], 10000.0))
        #     tens[r] = x[1]

        # if x[0] > 100000 and x[0] <= 1000000:
        #     # rount to nearest value of 100000
        #     r = int(roundup(x[0], 100000.0))
        #     tens[r] = x[1]

    for t, c in tens.iteritems():
        percent_log.append([t,c])

    percent_log.sort(cmp)

    other = []
    for t in percent_log:
        other.append([t[0],t[1]])

    with open("./results/filesPerProjectPercent_converted.json", "w") as text_file:
        text_file.write(json.dumps(other, sort_keys=True))


def roundup(x, size=10.0):
    import math
    return int(int(math.ceil(x / size)) * size)


def projectsPerLanguage():
    '''
    Converts output so we can see the number of projects per language
    https://jsfiddle.net/15v3c2pk/
    '''
    data = '[]'
    with open('./results/projectsPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x,y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1] == b[1]:
            return 0
        if a[1] > b[1]:
            return -1
        return 1

    new.sort(cmp)

    with open("./results/projectsPerLanguage_converted.json", "w") as text_file:
        text_file.write(json.dumps(new))


def mostCommonFileNames():
    '''
    Converts output so we can see the nmost common filenames
    '''
    data = '[]'
    with open('./results/fileNamesNoExtensionLowercaseCount.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    d = sorted(d.items(), key=operator.itemgetter(1), reverse=True)
    d = d[:51]

    with open("./results/fileNamesNoExtensionLowercaseCount_converted.json", "w") as text_file:
        text_file.write(json.dumps(d))


def largestPerLanguage():
    '''
    Convert the largest files per language into markdown
    table for embedding
    '''
    data = '[]'
    with open('./results/largestPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1]['Bytes'] == b[1]['Bytes']:
            return 0
        if a[1]['Bytes'] > b[1]['Bytes']:
            return -1
        return 1

    new.sort(cmp)

    res = [
        '| language | filename | bytes |',
        '| -------- | -------- | ----- |',
    ]

    for y in new:
        if 'bitbucket' in y[1]['Url']:
            link = y[1]['Url'] + '/src/master/' + '/'.join(y[1]['Location'].split('/')[3:])
        else:
            link = y[1]['Url'] + '/' + '/'.join(y[1]['Location'].split('/')[3:])

        x = '| %s | <a href="%s">%s</a> | %s |' % (y[0], link, y[1]['Filename'], y[1]['Value'])
        res.append(x)

    with open("./results/largestPerLanguage_converted.txt", "w") as text_file:
        text_file.write('''\n'''.join(res))


def longestPerLanguage():
    '''
    Convert the longest files per language into markdown
    table for embedding
    '''
    data = '[]'
    with open('./results/longestPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1]['Lines'] == b[1]['Lines']:
            return 0
        if a[1]['Lines'] > b[1]['Lines']:
            return -1
        return 1

    new.sort(cmp)

    res = [
        '| language | filename | lines |',
        '| -------- | -------- | ----- |',
    ]

    for y in new:
        if 'bitbucket' in y[1]['Url']:
            link = y[1]['Url'] + '/src/master/' + '/'.join(y[1]['Location'].split('/')[3:])
        else:
            link = y[1]['Url'] + '/' + '/'.join(y[1]['Location'].split('/')[3:])

        x = '| %s | <a href="%s">%s</a> | %s |' % (y[0], link, y[1]['Filename'], y[1]['Value'])
        res.append(x)

    with open("./results/longestPerLanguage_converted.txt", "w") as text_file:
        text_file.write('''\n'''.join(res))

def mostCommentedPerLanguage():
    '''
    Convert the longest files per language into markdown
    table for embedding
    '''
    data = '[]'
    with open('./results/mostCommentedPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1]['Comment'] == b[1]['Comment']:
            return 0
        if a[1]['Comment'] > b[1]['Comment']:
            return -1
        return 1

    new.sort(cmp)

    res = [
        '| language | filename | comment lines |',
        '| -------- | -------- | ------------- |',
    ]

    for y in new:
        if 'bitbucket' in y[1]['Url']:
            link = y[1]['Url'] + '/src/master/' + '/'.join(y[1]['Location'].split('/')[3:])
        else:
            link = y[1]['Url'] + '/' + '/'.join(y[1]['Location'].split('/')[3:])

        x = '| %s | <a href="%s">%s</a> | %s |' % (y[0], link, y[1]['Filename'], y[1]['Value'])
        res.append(x)

    with open("./results/mostCommentedPerLanguage_converted.txt", "w") as text_file:
        text_file.write('''\n'''.join(res))


def pureProjects():
    '''
    Converts the output of pureProjects into something
    we can throw into a chart library since it needs to 
    be sorted
    It is a count of the number of languages used by a project

    EG. languages:project where 123 projects have 2 languages in them
    https://jsfiddle.net/jqt81ufs/
    '''
    data = '[]'
    with open('./results/pureProjects.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([int(x), y])

    def cmp(a, b):
        if a == b:
            return 0
        if a < b:
            return -1
        return 1

    new.sort(cmp)

    with open("./results/pureProjects_converted.json", "w") as text_file:
        text_file.write(json.dumps(new, sort_keys=True))


def multipleGitIgnore():
    '''
    Converts the output of multipleGitIgnore into something
    we can throw into a chart library since it needs to 
    be sorted
    It is a count of the number of projects and gitignores

    EG. gitignore:project where 123 projects have 2 gitignore files
    https://jsfiddle.net/jqt81ufs/1/
    '''
    data = '[]'
    with open('./results/multipleGitIgnore.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([int(x), y])

    def cmp(a, b):
        if a == b:
            return 0
        if a < b:
            return -1
        return 1

    new.sort(cmp)

    with open("./results/multipleGitIgnore_converted.json", "w") as text_file:
        text_file.write(json.dumps(new, sort_keys=True))


def mostComplexPerLanguage():
    '''
    Convert the largest files per language into markdown
    table for embedding
    '''
    data = '[]'
    with open('./results/mostComplexPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1]['Value'] == b[1]['Value']:
            return 0
        if a[1]['Value'] > b[1]['Value']:
            return -1
        return 1

    new.sort(cmp)

    res = [
        '| language | filename | complexity |',
        '| -------- | -------- | ----- |',
    ]

    for y in new:
        if 'bitbucket' in y[1]['Url']:
            link = y[1]['Url'] + '/src/master/' + '/'.join(y[1]['Location'].split('/')[3:])
        else:
            link = y[1]['Url'] + '/' + '/'.join(y[1]['Location'].split('/')[3:])

        x = '| %s | <a href="%s">%s</a> | %s |' % (y[0], link, y[1]['Filename'], y[1]['Value'])
        res.append(x)

    with open("./results/mostComplexPerLanguage_converted.txt", "w") as text_file:
        text_file.write('''\n'''.join(res))


if __name__ == '__main__':
    filesPerProject()
    filesPerProjectPercentile()
    projectsPerLanguage()
    mostCommonFileNames()
    largestPerLanguage()
    longestPerLanguage()
    pureProjects()
    multipleGitIgnore()
    mostComplexPerLanguage()
    mostCommentedPerLanguage()



'''
files
java 100
php 50
c 150

complexity
java 300
php 200
c 500


100 / 150 = 0.67
50 / 150 = 0.34
150 / 150 = 1


weighted complexity


java = 201
php = 68
c = 500
'''

'''
average complexity of Java repo
average complexity of Java repo between 1-50 files
average complexity of Java repo between 51-100 files
etc...
'''