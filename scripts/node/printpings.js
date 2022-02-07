import axios from "axios";

const hosts = [
    "https://yumseng-m4pgrqojya-as.a.run.app",
    "https://yumseng.herokuapp.com",
    "https://yumseng-nnwor.ondigitalocean.app"
]

async function main() {
    console.log(JSON.stringify(await Promise.all(hosts.map(
        async (host) => {
            const resp = await axios.get(`${host}/ping`)
            const data = resp.data
            data["host"] = host
            data["readable_time"] = new Intl.DateTimeFormat('en-SG', {
                dateStyle: 'full',
                timeStyle: 'long'
            }).format(new Date(data["server_start_time"]))
            return data
        })
    ), null, '\t'))
}

main()

