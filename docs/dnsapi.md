
## 1. Use CloudFlare domain API to automatically issue cert
First you need to login to your CloudFlare account to get your [API key](https://dash.cloudflare.com/profile). 
```
.spec.env.CF_Key: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.env.CF_Email: "xxxx@sss.com"
.spec.tpye: dns_cf
```

## 2. Use DNSPod.cn domain API to automatically issue cert
First you need to login to your DNSPod account to get your API Key and ID.
```
.spec.env.DP_Id: "1234"
.spec.env.DP_Key: "sADDsdasdgdsf"
.spec.tpye: dns_dp
```

## 3. Use CloudXNS.com domain API to automatically issue cert
First you need to login to your CloudXNS account to get your API Key and Secret.
```
.spec.env.CX_Key: "1234"
.spec.env.CX_Secret: "sADDsdasdgdsf"
.spec.tpye: dns_cx
```

## 4. Use GoDaddy.com domain API to automatically issue cert
First you need to login to your GoDaddy account to get your API Key and Secret.
https://developer.godaddy.com/keys/
```
.spec.env.GD_Key: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.env.GD_Secret: "asdfsdafdsfdsfdsfdsfdsafd"
.spec.tpye: dns_gd
```

## 5. Use PowerDNS embedded API to automatically issue cert
First you need to login to your PowerDNS account to enable the API and set your API-Token in the configuration.
https://doc.powerdns.com/md/httpapi/README/
```
.spec.env.PDNS_Url: "http://ns.example.com:8081"
.spec.env.PDNS_ServerId: "localhost"
.spec.env.PDNS_Token: "0123456789ABCDEF"
.spec.env.PDNS_Ttl: 60
.spec.tpye: dns_pdns
```

## 6. Use OVH/kimsufi/soyoustart/runabove API to automatically issue cert



## 7. Use nsupdate to automatically issue cert
First, generate a key for updating the zone
```
b=$(dnssec-keygen -a hmac-sha512 -b 512 -n USER -K /tmp foo)
cat > /etc/named/keys/update.key <<EOF
key "update" {
    algorithm hmac-sha512;
    secret "$(awk '/^Key/{print $2}' /tmp/$b.private)";
};
EOF
rm -f /tmp/$b.{private,key}
```

Include this key in your named configuration
```
include "/etc/named/keys/update.key";
```

Next, configure your zone to allow dynamic updates.

Depending on your named version, use either
```
zone "example.com" {
    type master;
    allow-update { key "update"; };
};
```
or
```
zone "example.com" {
    type master;
    update-policy {
        grant update subdomain example.com.;
    };
}
```

Finally


```
.spec.env.NSUPDATE_SERVER: "dns.example.com"
.spec.env.NSUPDATE_KEY: "/path/to/your/nsupdate.key"
.spec.env.NSUPDATE_ZONE: "example.com"
.spec.tpye: dns_nsupdate
```

## 8. Use LuaDNS domain API
Get your API token at https://api.luadns.com/settings
```
.spec.env.LUA_Key: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.env.LUA_Email: "xxxx@sss.com"
.spec.tpye: dns_lua
```
## 9. Use DNSMadeEasy domain API
Get your API credentials at https://cp.dnsmadeeasy.com/account/info
```
.spec.env.ME_Key: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.env.ME_Secret: "qdfqsdfkjdskfj"
.spec.tpye: dns_me
```

## 10. Use Amazon Route53 domain API
```
.spec.env. AWS_ACCESS_KEY_ID: XXXXXXXXXX
.spec.env. AWS_SECRET_ACCESS_KEY: XXXXXXXXXXXXXXX
.spec.tpye: dns_aws
```

## 11. Use Aliyun domain API to automatically issue cert
First you need to login to your Aliyun account to get your API key.
[https://ak-console.aliyun.com/#/accesskey](https://ak-console.aliyun.com/#/accesskey)
```
.spec.env.Ali_Key: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.env.Ali_Secret: "jlsdflanljkljlfdsaklkjflsa"
.spec.tpye: dns_ali
```

## 12. Use ISPConfig 3.1 API
This only works for ISPConfig 3.1 (and newer).

Create a Remote User in the ISPConfig Control Panel. The Remote User must have access to at least `DNS zone functions` and `DNS txt functions`.


```
.spec.env.ISPC_User: "xxx"
.spec.env.ISPC_Password: "xxx"
.spec.env.ISPC_Api: "https://ispc.domain.tld:8080/remote/json.php"
.spec.env.ISPC_Api_Insecure: 1
.spec.tpye: dns_ispconfig
```

## 13. Use Alwaysdata domain API
First you need to login to your Alwaysdata account to get your API Key.
```
.spec.env.AD_API_KEY: "myalwaysdataapikey"
.spec.tpye: dns_ad
```

## 14. Use Linode domain API


### Classic Manager ###
Classic Manager: https://manager.linode.com/profile/api
First you need to login to your Linode account to get your API Key.

Then add an API key with label *ACME* and copy the new key into the following
command.


```
.spec.env.LINODE_API_KEY: "..."
.spec.tpye: dns_linode
```

## 15. Use FreeDNS
FreeDNS (https://freedns.afraid.org/) does not provide an API to update DNS records (other than IPv4 and IPv6


```
.spec.env.FREEDNS_User: "..."
.spec.env.FREEDNS_Password: "..."
.spec.tpye: dns_freedns
```

## 16. Use cyon.ch
```
.spec.env.CY_Username: "your_cyon_username"
.spec.env.CY_Password: "your_cyon_password"
.spec.env.CY_OTP_Secret: "your_otp_secret" # Only required if using 2FA
.spec.tpye: dns_cyon
```

## 17. Use Domain-Offensive/Resellerinterface/Domainrobot API
You will need your login credentials (Partner ID+Password) to the Resellerinterface
```
.spec.env.DO_PID: "KD-1234567"
.spec.env.DO_PW: "cdfkjl3n2"
.spec.tpye: dns_do
```

## 18. Use Gandi LiveDNS API
You must enable the new Gandi LiveDNS API first and the create your api key, See: http://doc.livedns.gandi.net/
```
.spec.env.GANDI_LIVEDNS_KEY: "fdmlfsdklmfdkmqsdfk"
.spec.tpye: dns_gandi_livedns
```

## 19. Use Knot (knsupdate) DNS API to automatically issue cert
First, generate a TSIG key for updating the zone.

```
keymgr tsig generate -t acme_key hmac-sha512 > /etc/knot/acme.key
```

Include this key in your knot configuration file.

```
include: /etc/knot/acme.key
```

Next, configure your zone to allow dynamic updates.

Dynamic updates for the zone are allowed via proper ACL rule with the `update` action. For in-depth instructions, please see [Knot DNS's documentation](https://www.knot-dns.cz/documentation/).

```
acl:
  - id: acme_acl
    address: 192.168.1.0/24
    key: acme_key
    action: update

zone:
  - domain: example.com
    file: example.com.zone
    acl: acme_acl
```

Finally


```
.spec.env.KNOT_SERVER: "dns.example.com"
.spec.env.KNOT_KEY: `grep \# /etc/knot/acme.key | cut -d' ' -f2`
.spec.tpye: dns_knot
```

## 20. Use DigitalOcean API (native)
You need to obtain a read and write capable API key from your DigitalOcean account. See: https://www.digitalocean.com/help/api/
```
.spec.env.DO_API_KEY: "75310dc4ca779ac39a19f6355db573b49ce92ae126553ebd61ac3a3ae34834cc"
.spec.tpye: dns_dgon
```

## 21. Use ClouDNS.net API
You need to set the HTTP API user ID and password credentials. See: https://www.cloudns.net/wiki/article/42/. For security reasons, it's recommended to use a sub user ID that only has access to the necessary zones, as a regular API user has access to your entire account.
```
.spec.env.CLOUDNS_SUB_AUTH_ID: XXXXX
.spec.env.CLOUDNS_AUTH_PASSWORD: "YYYYYYYYY"
.spec.tpye: dns_cloudns
```

## 22. Use Infoblox API
First you need to create/obtain API credentials on your Infoblox appliance.
```
.spec.env.Infoblox_Creds: "username:password"
.spec.env.Infoblox_Server: "ip or fqdn of infoblox appliance"
.spec.tpye: dns_infoblox
```

## 23. Use VSCALE API
First you need to create/obtain API tokens on your [settings panel](https://vscale.io/panel/settings/tokens/).
```
.spec.env.VSCALE_API_KEY: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.tpye: dns_vscale
```

##  24. Use Dynu API
First you need to create/obtain API credentials from your Dynu account. See: https://www.dynu.com/resources/api/documentation
```
.spec.env.Dynu_ClientId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
.spec.env.Dynu_Secret: "yyyyyyyyyyyyyyyyyyyyyyyyy"
.spec.tpye: dns_dynu
```

## 25. Use DNSimple API
First you need to login to your DNSimple account and generate a new oauth token.
https://dnsimple.com/a/{your account id}/account/access_tokens
```
.spec.env.DNSimple_OAUTH_TOKEN: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.tpye: dns_dnsimple
```

## 26. Use NS1.com API
```
.spec.env.NS1_Key: "fdmlfsdklmfdkmqsdfk"
.spec.tpye: dns_nsone
```

## 27. Use DuckDNS.org API



## 28. Use Name.com API
Create your API token here: https://www.name.com/account/settings/api
```
.spec.env.Namecom_Username: "testuser"
.spec.env.Namecom_Token: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
.spec.tpye: dns_namecom
```

## 29. Use Dyn Managed DNS API to automatically issue cert
First, login to your Dyn Managed DNS account: https://portal.dynect.net/login/
```
.spec.env.DYN_Customer: "customer"
.spec.env.DYN_Username: "apiuser"
.spec.env.DYN_Password: "secret"
.spec.tpye: dns_dyn
```

## 30. Use pdd.yandex.ru API
```
.spec.env.PDD_Token: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
.spec.tpye: dns_yandex
```

## 31. Use Hurricane Electric
Hurricane Electric (https://dns.he.net/) doesn't have an API so just set your login credentials like so:
```
.spec.env.HE_Username: "yourusername"
.spec.env.HE_Password: "password"
.spec.tpye: dns_he
```

## 32. Use UnoEuro API to automatically issue cert
First you need to login to your UnoEuro account to get your API key.
```
.spec.env.UNO_Key: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.env.UNO_User: "UExxxxxx"
.spec.tpye: dns_unoeuro
```

## 33. Use INWX
[INWX](https://www.inwx.de/) offers an [xmlrpc api](https://www.inwx.de/de/help/apidoc)  with your standard login credentials, set them like so:
```
.spec.env.INWX_User: "yourusername"
.spec.env.INWX_Password: "password"

.spec.env.INWX_Shared_Secret: "shared secret"

.spec.tpye: dns_inwx
```


## 34. User Servercow API v1
```
.spec.env.SERVERCOW_API_Username: username
.spec.env.SERVERCOW_API_Password: password
.spec.tpye: dns_servercow
```

## 35. Use Namesilo.com API
You'll need to generate an API key at https://www.namesilo.com/account_api.php
```
.spec.env.Namesilo_Key: "xxxxxxxxxxxxxxxxxxxxxxxx"
.spec.tpye: dns_namesilo
```

## 36. Use autoDNS (InternetX)
[InternetX](https://www.internetx.com/) offers an [xml api](https://help.internetx.com/display/API/AutoDNS+XML-API)  with your standard login credentials, set them like so:
```
.spec.env.AUTODNS_USER: "yourusername"
.spec.env.AUTODNS_PASSWORD: "password"
.spec.env.AUTODNS_CONTEXT: "context"
.spec.tpye: dns_autodns
```

## 37. Use Azure DNS
You have to create a service principal first. See:[How to use :[How to use Azure DNS](../../../wiki/How-to-use-Azure-DNS)


```
.spec.env.AZUREDNS_SUBSCRIPTIONID: "12345678-9abc-def0-1234-567890abcdef"
.spec.env.AZUREDNS_TENANTID: "11111111-2222-3333-4444-555555555555"
.spec.env.AZUREDNS_APPID: "3b5033b5-7a66-43a5-b3b9-a36b9e7c25ed"
.spec.env.AZUREDNS_CLIENTSECRET: "1b0224ef-34d4-5af9-110f-77f527d561bd"
.spec.tpye: dns_azure
```

## 38. Use selectel.com(selectel.ru) domain API to automatically issue cert
First you need to login to your account to get your API key from: https://my.selectel.ru/profile/apikeys.
```
.spec.env.SL_Key: "sdfsdfsdfljlbjkljlkjsdfoiwje"
.spec.tpye: dns_selectel
```

## 39. Use zonomi.com domain API to automatically issue cert
First you need to login to your account to find your API key from: http://zonomi.com/app/dns/dyndns.jsp

Your will find your api key in the example urls:

```sh
https://zonomi.com/app/dns/dyndns.jsp?host=example.com&api_key=1063364558943540954358668888888888
```
```
.spec.env.ZM_Key: "1063364558943540954358668888888888"
.spec.tpye: dns_zonomi
```


## 40. Use DreamHost DNS API
DNS API keys may be created at https://panel.dreamhost.com/?tree: home.api.

Ensure the created key has add and remove privelages.


```

.spec.env.DH_API_KEY: "<api key>"
.spec.tpye: dns_dreamhost
```

## 41. Use DirectAdmin API
See https://www.directadmin.com/api.php and https://www.directadmin.com/features.php?id: 1298


```
.spec.env.DA_Api: "https://remoteUser:remotePassword@da.domain.tld:8443"
.spec.env.DA_Api_Insecure: 1
.spec.tpye: dns_da
```

## 42. Use KingHost DNS API
API access must be enabled at https://painel.kinghost.com.br/painel.api.php


```
.spec.env.KINGHOST_Username: "yourusername"
.spec.env.KINGHOST_Password: "yourpassword"
.spec.tpye: dns_kinghost
```

## 43. Use Zilore DNS API
First, get your API key at https://my.zilore.com/account/api


```
.spec.env.Zilore_Key: "5dcad3a2-36cb-50e8-cb92-000002f9"
.spec.tpye: dns_zilore
```

## 44. Use Loopia.se API

```
.spec.env.LOOPIA_User: "user@loopiaapi"
.spec.env.LOOPIA_Password: "password"
.spec.tpye: dns_loopia
```

## 45. Use ACME DNS API
https://github.com/joohoi/acme-dns


```
.spec.env.ACMEDNS_UPDATE_URL: "https://auth.acme-dns.io/update"
.spec.env.ACMEDNS_USERNAME: "<username>"
.spec.env.ACMEDNS_PASSWORD: "<password>"
.spec.env.ACMEDNS_SUBDOMAIN: "<subdomain>"
.spec.tpye: dns_acmedns
```

## 46. Use TELE3 API
First you need to login to your TELE3 account to set your API-KEY.
https://www.tele3.cz/system-acme-api.html

```

.spec.env.TELE3_Key: "MS2I4uPPaI..."
.spec.env.TELE3_Secret: "kjhOIHGJKHg"
.spec.tpye: dns_tele3
```

## 47. Use Euserv.eu API
First you need to login to your euserv.eu account and activate your API Administration (API Verwaltung).
[https://support.euserv.com](https://support.euserv.com)

Once you've activate, login to your API Admin Interface and create an API account.
Please specify the scope (active groups: domain) and assign the allowed IPs.

```

.spec.env.EUSERV_Username: "99999.user123"
.spec.env.EUSERV_Password: "Asbe54gHde"
.spec.tpye: dns_euserv
```


## 48. Use DNSPod.com domain API to automatically issue cert
First you need to get your API Key and ID by this [get-the-user-token](https://www.dnspod.com/docs/info.html#get-the-user-token).


```
.spec.env.DPI_Id: "1234"
.spec.env.DPI_Key: "sADDsdasdgdsf"
.spec.tpye: dns_dpi
```

## 49. Use Google Cloud DNS API to automatically issue cert
First you need to authenticate to gcloud.

```
gcloud init
```

**The `dns_gcloud` script uses the active gcloud configuration and credentials.**
There is no logic inside `dns_gcloud` to override the project and other settings.
If needed, create additional [gcloud configurations](https://cloud.google.com/sdk/gcloud/reference/topic/configurations).
You can change the configuration being used without *activating* it; simply set the `CLOUDSDK_ACTIVE_CONFIG_NAME` environment variable.

```

.spec.env.CLOUDSDK_ACTIVE_CONFIG_NAME: default  # see the note above
.spec.tpye: dns_gcloud
```

## 50. Use ConoHa API
First you need to login to your ConoHa account to get your API credentials.


```
.spec.env.CONOHA_Username: "xxxxxx"
.spec.env.CONOHA_Password: "xxxxxx"
.spec.env.CONOHA_TenantId: "xxxxxx"
.spec.env.CONOHA_IdentityServiceApi: "https://identity.xxxx.conoha.io/v2.0"
.spec.tpye: dns_conoha
```

## 51. Use netcup DNS API to automatically issue cert
First you need to login in your CCP account to get your API Key and API Password.


```
.spec.env.NC_Apikey: "<Apikey>"
.spec.env.NC_Apipw: "<Apipassword>"
.spec.env.NC_CID: "<Customernumber>"
.spec.tpye: dns_netcup
```

## 52. Use GratisDNS.dk

## 53. Use Namecheap
You will need your namecheap username, API KEY (https://www.namecheap.com/support/api/intro.aspx) and your external IP address (or an URL to get it), this IP will need to be whitelisted at Namecheap.


```
.spec.env.NAMECHEAP_USERNAME: "..."
.spec.env.NAMECHEAP_API_KEY: "..."
.spec.env.NAMECHEAP_SOURCEIP: "..."
NAMECHEAP_SOURCEIP can either be an IP address or an URL to provide it (e.g. https://ifconfig.co/ip).
.spec.tpye: dns_namecheap
```

## 54. Use MyDNS.JP API
First, register to MyDNS.JP and get MasterID and Password.

```

.spec.env.MYDNSJP_MasterID: MasterID
.spec.env.MYDNSJP_Password: Password
.spec.tpye: dns_mydnsjp
```

## 55. Use hosting.de API
Create an API key in your hosting.de account here: https://secure.hosting.de

```

.spec.env.HOSTINGDE_APIKEY: 'xxx'
.spec.env.HOSTINGDE_ENDPOINT: 'https://secure.hosting.de'
.spec.tpye: dns_hostingde
```

## 56. Use Neodigit.net API
```
.spec.env.NEODIGIT_API_TOKEN: "eXJxTkdUVUZmcHQ3QWJackQ4ZGlMejRDSklRYmo5VG5zcFFKK2thYnE0WnVnNnMy"
.spec.tpye: dns_neodigit
```

## 57. Use Exoscale API
```
.spec.env.EXOSCALE_API_KEY: 'xxx'
.spec.env.EXOSCALE_SECRET_KEY: 'xxx'
.spec.tpye: dns_exoscale
```

## 58. Using PointHQ API to issue certs
Log into [PointHQ account management](https://app.pointhq.com/profile) and copy the API key from the page there.


```
.spec.env.PointHQ_Key: "apikeystringgoeshere"
.spec.env.PointHQ_Email: "accountemail@yourdomain.com"
.spec.tpye: dns_pointhq
```

## 59. Use Active24 API
Create an API token in the Active24 account section, documentation on https://faq.active24.com/cz/790131-REST-API-rozhran%C3%AD.

```

.spec.env.ACTIVE24_Token: 'xxx'
.spec.tpye: dns_active24

```
## 60. Use do.de API
Create an API token in your do.de account ([Create token here](https://www.do.de/account/letsencrypt/) | [Documentation](https://www.do.de/wiki/LetsEncrypt_-_Entwickler)).

```

.spec.env.DO_LETOKEN: 'FmD408PdqT1E269gUK57'
.spec.tpye: dns_doapi
```

## 61. Use Nexcess API



## 62. Use Thermo.io API



## 63. Use Futurehosting API



## 64. Use Rackspace API
```
.spec.env.RACKSPACE_Username: 'username'
.spec.env.RACKSPACE_Apikey: 'xxx'
.spec.tpye: dns_rackspace
```

## 65. Use Online API
First, you'll need to retrive your API key, which is available under https://console.online.net/en/api/access


```
.spec.env.ONLINE_API_KEY: 'xxx'
.spec.tpye: dns_online
```

## 66. Use MyDevil.net



## 67. Use Core-Networks API to automatically issue cert
First you need to login to your Core-Networks account to to set up an API-User.
Then .spec.env.username and password to use these credentials.


```
.spec.env.CN_User: "user"
.spec.env.CN_Password: "passowrd"
.spec.tpye: dns_cn
```

## 68. Use NederHost API
```
.spec.env.NederHost_Key: 'xxx'
.spec.tpye: dns_nederhost
```

## 69. Use Zone.ee DNS API
First, you'll need to retrive your API key. Estonian insructions https://help.zone.eu/kb/zoneid-api-v2/


```
.spec.env.ZONE_Username: yourusername
.spec.env.ZONE_Key: keygoeshere
acme.sh --issue -d example.com -d www.example.com --dns dns_zone
```
## 70. Use UltraDNS API

UltraDNS is a paid for service that provides DNS, as well as Web and Mail forwarding (as well as reporting, auditing, and advanced tools).

More information can be found here: https://www.security.neustar/lp/ultra20/index.html

The REST API documentation for this service is found here: https://portal.ultradns.com/static/docs/REST-API_User_Guide.pdf 

Set your UltraDNS User name, and password; these would be the same you would use here:

https://portal.ultradns.com/ - or if you create an API only user, that username and password would be better utilized.

```

.spec.env.ULTRA_USR: "abcd"
.spec.env.ULTRA_PWD: "efgh"
.spec.tpye: dns_ultra
```

## 71. Use deSEC.io
Sign up for dynDNS at https://desec.io first.

Set your API token (password) and domain (username) from your email sent by desec.io


```
.spec.env.DEDYN_TOKEN: d41d8cd98f00b204e9800998ecf8427e
.spec.env.DEDYN_NAME: foobar.dedyn.io
.spec.tpye: dns_desec
```

## 72. Use OpenProvider API
First, you need to enable API access and retrieve your password hash on https://rcp.openprovider.eu/account/dashboard.php

```

.spec.env.OPENPROVIDER_USER: 'username'
.spec.env.OPENPROVIDER_PASSWORDHASH: 'xxx'
.spec.tpye: dns_openprovider
```

## 73. Use MaraDNS API
Make sure you've configured MaraDNS properly and setup a zone file for your domain. See [`csv2(5)`](https://manpages.debian.org/stretch/maradns/csv2.5.en.html).
Set the path to your zone file, and path to duende's pid file (see, [`duende(8)`](https://manpages.debian.org/stretch/duende/duende.8.en.html) or `ps -C duende o pid,cmd`).

```

.spec.env.MARA_ZONE_FILE: "/etc/maradns/db.domain.com"
.spec.env.MARA_DUENDE_PID_PATH: "/run/maradns/etc_maradns_mararc.pid"
.spec.tpye: dns_maradns
.spec.tpye: dns_myapi
```


