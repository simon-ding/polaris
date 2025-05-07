package nat

import (
	"fmt"
	"strings"
	"polaris/log"

	"github.com/pion/stun/v3"
)

func getNatIpAndPort() (*stun.XORMappedAddress, error) {

	var xorAddr stun.XORMappedAddress

	for _, server := range getStunServers() {
		log.Infof("try to connect to stun server: %s", server)
		u, err := stun.ParseURI("stun:" + server)
		if err != nil {
			continue
		}
		// Creating a "connection" to STUN server.
		c, err := stun.DialURI(u, &stun.DialConfig{})
		if err != nil {
			continue
		}
		// Building binding request with random transaction id.
		message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
		// Sending request to STUN server, waiting for response message.
		var err1 error
		if err := c.Do(message, func(res stun.Event) {
			if res.Error != nil {
				err1 = res.Error
				return
			}
			log.Infof("stun server %s response: %v", server, res.Message.String())
			// Decoding XOR-MAPPED-ADDRESS attribute from message.

			if err := xorAddr.GetFrom(res.Message); err != nil {
				err1 = err
				return
			}
			fmt.Println("your IP is", xorAddr.IP)
			fmt.Println("your port is", xorAddr.Port)
		}); err != nil {
			log.Warnf("stun server %s error: %v", server, err)
			continue
		}
		if err1 != nil {
			log.Warnf("stun server %s error: %v", server, err1)
			continue
		}
		break
	}
	return &xorAddr, nil
}

func getStunServers() []string {
	var servers []string
	for _, line := range strings.Split(strings.TrimSpace(stunServers1), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		servers = append(servers, line)
	}
	return servers
}

// https://github.com/heiher/natmap/issues/18
const stunServers1 = `
stun.miwifi.com:3478
stun.chat.bilibili.com:3478
stun.cloudflare.com:3478
turn.cloudflare.com:3478
fwa.lifesizecloud.com:3478
`

// https://github.com/pradt2/always-online-stun
const stunServers = `
stun.miwifi.com:3478
stun.ukh.de:3478
stun.kanojo.de:3478
stun.m-online.net:3478
stun.nextcloud.com:3478
stun.voztovoice.org:3478
stun.oncloud7.ch:3478
stun.antisip.com:3478
stun.bitburger.de:3478
stun.acronis.com:3478
stun.signalwire.com:3478
stun.sonetel.net:3478
stun.poetamatusel.org:3478
stun.avigora.fr:3478
stun.diallog.com:3478
stun.nanocosmos.de:3478
stun.romaaeterna.nl:3478
stun.heeds.eu:3478
stun.freeswitch.org:3478
stun.engineeredarts.co.uk:3478
stun.root-1.de:3478
stun.healthtap.com:3478
stun.allflac.com:3478
stun.vavadating.com:3478
stun.godatenow.com:3478
stun.mixvoip.com:3478
stun.sip.us:3478
stun.sipthor.net:3478
stun.stochastix.de:3478
stun.kaseya.com:3478
stun.files.fm:3478
stun.meetwife.com:3478
stun.myspeciality.com:3478
stun.3wayint.com:3478
stun.voip.blackberry.com:3478
stun.axialys.net:3478
stun.bridesbay.com:3478
stun.threema.ch:3478
stun.siptrunk.com:3478
stun.ncic.com:3478
stun.1cbit.ru:3478
stun.ttmath.org:3478
stun.yesdates.com:3478
stun.sonetel.com:3478
stun.peethultra.be:3478
stun.pure-ip.com:3478
stun.business-isp.nl:3478
stun.ringostat.com:3478
stun.imp.ch:3478
stun.cope.es:3478
stun.baltmannsweiler.de:3478
stun.lovense.com:3478
stun.frozenmountain.com:3478
stun.linuxtrent.it:3478
stun.thinkrosystem.com:3478
stun.3deluxe.de:3478
stun.skydrone.aero:3478
stun.ru-brides.com:3478
stun.streamnow.ch:3478
stun.atagverwarming.nl:3478
stun.ipfire.org:3478
stun.fmo.de:3478
stun.moonlight-stream.org:3478
stun.f.haeder.net:3478
stun.nextcloud.com:443
stun.finsterwalder.com:3478
stun.voipia.net:3478
stun.zepter.ru:3478
stun.sipnet.net:3478
stun.hot-chilli.net:3478
stun.zentauron.de:3478
stun.geesthacht.de:3478
stun.annatel.net:3478
stun.flashdance.cx:3478
stun.voipgate.com:3478
stun.genymotion.com:3478
stun.graftlab.com:3478
stun.fitauto.ru:3478
stun.telnyx.com:3478
stun.verbo.be:3478
stun.dcalling.de:3478
stun.lleida.net:3478
stun.romancecompass.com:3478
stun.siplogin.de:3478
stun.bethesda.net:3478
stun.alpirsbacher.de:3478
stun.uabrides.com:3478
stun.technosens.fr:3478
stun.radiojar.com:3478
`
