const host = "https://yumseng-nnwor.ondigitalocean.app"

import axios from 'axios'
async function seng(userId, roomId, imageUrl = "") {
    console.log(` ${userId} sending cheer ${roomId}`)
    await axios.post(`${host}/rooms/${roomId}/user/${userId}/cheers`, {
        value: "seng",
        client_created_at: new Date().toJSON(),
        user_id: userId,
        image_url: imageUrl
    })
    console.log(`${userId} sent cheer ${roomId}`)
}

async function enableSeng(roomId) {
    await axios.delete(`${host}/rooms/${roomId}/cheer-rules`)
}

(async function () {
    await enableSeng('global-room')
    await seng("POST_OFFIC", "global-room")
})()