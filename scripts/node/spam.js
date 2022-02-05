import axios from 'axios'
import testToImage from 'text-to-image'

const host = "https://yumseng-m4pgrqojya-as.a.run.app/" // "https://yumseng.herokuapp.com"

async function joinRoom(userId, roomId) {
    console.log(`user ${userId} joining room ${roomId}`)
    const response = await axios.post(`${host}/rooms/${roomId}/user/${userId}`)
    console.log(`user ${userId} joined room ${roomId} ${response.data}`)
}

async function pingServer() {
    return await axios.get(`${host}/ping`)
}

async function yum(userId, roomId, imageUrl = "") {
    console.log(` ${userId} sending cheer ${roomId}`)
    await axios.post(`${host}/rooms/${roomId}/user/${userId}/cheers`, {
        value: "yum",
        client_created_at: new Date().toJSON(),
        user_id: userId,
        image_url: imageUrl
    })
    console.log(`${userId} sent cheer ${roomId}`)
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function main() {
    console.log("pinging server")

    const pingResponse = await pingServer()

    console.log(pingResponse.data)

    const userIds = [
        "Aaren",
        "Aarika",
        "Abagael",
        "Abagail",
        "Abbe",
        "Abbey",
        "Abbi",
        "Abbie",
        "Abby",
        "Abbye",
        "Abigael",
        "Abigail",
        "Abigale",
        "Abra",
        "Ada",
        "Adah",
        "Adaline",
        "Adan",
        "Adara",
        "Adda",
        "Addi",
        "Addia",
        "Addie",
        "Addy",
        "Adel",
        "Adela",
        "Adelaida",
        "Adelaide",
        "Adele",
        "Adelheid",
        "Adelice",
        "Adelina",
        "Adelind",
        "Adeline",
        "Adella",
        "Adelle",
        "Adena",
        "Adey",
        "Adi",
        "Adiana",
        "Adina",
        "Adora",
        "Adore",
        "Adoree",
        "Adorne",
        "Adrea",
        "Adria",
        "Adriaens",
        "Adrian",
        "Adriana",
        "Adriane",
        "Adrianna",
        "Adrianne"]

    async function generateImageUrl(userId, bgColor) {
        const config = {
            fontSize: 80,
            bgColor: bgColor,
            customHeight: 400
        }
        console.log(`image url config ${config}`)
        return await testToImage.generate(userId, config)
    }

    const bgColors = [
        '#FF0000',
        '#800000',
        '#FFFF00',
        '#808000',
        '#00FF00',
        '#008000',
        `#00FFFF`,
        `#008080`
    ]

    async function joinAndSend(userId, index) {
        try {
            userId = `${index}-${userId}`
            await joinRoom(userId, 'global-room')
            const imgUrl = await generateImageUrl(userId, bgColors[index % bgColors.length])
            for (let i = 0; i < 500; i++) {
                console.log(`${i}/200 ${userId}`)
                try {
                    await yum(userId, 'global-room', imgUrl)
                    await sleep(Math.random() * 50);
                } catch (e) {
                    console.log(`${i}/200 EXCEPTION ${userId} ${e}`)
                }
            }
        } catch (exception) {
            console.log(`${userId} ${exception}`)
        }
    }

    const all = userIds.map(joinAndSend)
    console.log("awaiting all jobs")
    await Promise.all(all)
    console.log("finished")
}

main()
