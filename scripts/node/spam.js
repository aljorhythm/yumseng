import axios from 'axios'
import testToImage from 'text-to-image'

// const host = "https://yumseng-m4pgrqojya-as.a.run.app" //
// const host = "http://localhost"
// "https://yumseng.herokuapp.com"
const host = "https://yumseng-nnwor.ondigitalocean.app"
const USERS_COUNT = 20

const CHEERS_PER_USER = 40000;
const TIMEOUT_SECONDS = 160;
const SLEEP_TIME = 200;

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

    const userIds = ["Ice_cream_cone_Post_office", "Clock_Plus", "Dog_Sink", "Crab_System", "Settings_Bird", "Urine_Running", "Floppy_Disk_Toolbox", "Robot_Shoe", "Ring_Male", "Crab_Floppy_Disk", "Plants_Kitty", "Shoe_Fence", "Rollers_Horse", "Cone_Allergies", "Running_Breakfast", "Fusion_Settings", "Male_Soap", "Poop_Flowers", "Solar_Plus", "Nuclear_Video_games", "Android_Cat", "Toilet_Boat", "Puppy_Plus", "Kitty_Cone", "Floppy_Disk_Light_saber", "Body_Toolbox", "Boat_Plus", "Toilet_Ice_cream_cone", "Whale_Comics", "Sink_Book", "Mail_Post_office", "Websites_Breakfast", "Puppy_Settings", "Websites_Breakfast", "Post_office_BBQ", "Fusion_Water", "Dislike_Horse", "Settings_Clock", "Water_Ice_cream", "Urine_Post_office", "Water_Boat", "Elevator_Ring", "Comics_Website", "Video_games_Shower", "Clock_Horse", "Shoes_Trees", "Soda_Shoes", "Robot_YouTube", "Settings_Mail", "Running_Toolbox", "Horse_Crab", "Clock_Plus", "Trees_Allergies", "Puppy_Plus", "Robot_Leash", "Prints_BBQ", "Cone_Shoe", "Toolbox_Solar", "Rollers_System", "Trees_Mail", "Website_Prints", "Bird_Allergies", "Plants_Towel", "Fence_Post_office", "Ice_cream_Allergies", "Soap_Male", "Toilet_Websites", "Plants_Whale", "Shoes_Rollers", "Whale_Elevator", "Prints_Whale", "Ice_cream_Elevator", "Video_games_Poop", "Toolbox_Cone", "Plants_Fusion", "Breakfast_Body", "Settings_Nuclear", "Websites_Shoes", "Hnads_Fence", "BBQ_Light_saber", "Clock_Post_office", "Shoe_Cone", "Ice_cream_cone_Video_games", "Settings_Comics", "BBQ_Prints", "Solar_Boat", "Android_Allergies", "Drugs_Dislike", "Printer_Urine", "Printer_Comics", "Horse_Shoe", "Breakfast_Breakfast", "YouTube_Toolbox", "Post_office_Allergies", "Fusion_Boat", "Solar_Light_saber", "Elevator_Breakfast", "Mail_Toolbox", "Ring_Leash", "Laptop_Floppy_Disk", "Toolbox_Comics", "Poop_Bird", "Post_office_Mail", "Post_office_Flowers", "Solar_Soda", "Soda_Comics", "Clock_Toolbox", "Shelf_Cone", "Soda_Light_saber", "Ice_cream_cone_Breakfast", "Urine_Android", "Dislike_Towel", "Cone_Soap", "Toilet_Dislike", "Toilet_BBQ", "Book_Shoes", "Fusion_Ring", "Ice_cream_cone_Nuclear", "Clock_Plus", "Urine_YouTube", "Cat_Website", "Video_games_Prints", "Flowers_Soap", "Laptop_Nuclear", "Prints_Laptop", "Mail_Water", "Plants_Poop", "Website_Male", "Hnads_Leash", "Android_Water", "Body_Toolbox", "Plants_Allergies", "Kitty_Mail", "Laptop_BBQ", "Prints_Light_saber", "Android_Dislike", "BBQ_Light_saber", "Shower_Soap", "Dislike_Whale", "Prints_Puppy", "Male_Bird", "Dislike_Websites", "Rollers_Bird", "Settings_Ice_cream_cone", "Puppy_Fence", "Post_office_Flowers", "Running_Cat", "Boat_Settings", "Kitty_Laptop", "Whale_Running", "Printer_Boat", "Urine_Running", "Nuclear_Water", "Cat_Soap", "Drugs_Video_games", "Flowers_Website", "Horse_Post_office", "Ice_cream_Kitty", "Boat_Leash", "Shoes_Clock", "Toolbox_Websites", "Drugs_Settings", "Flowers_Toilet", "Bird_Mail", "System_Plants", "Website_Dog", "Shelf_Settings", "Settings_Solar", "Dog_Kitty", "Post_office_Shoes", "Post_office_Comics", "Ice_cream_Flowers", "Websites_Post_office", "Poop_Boat", "Shelf_Prints", "Soap_Toolbox", "Floppy_Disk_Soap", "Toilet_YouTube", "Whale_Water", "Shelf_Dislike", "Ring_Shoes", "Urine_Prints", "Floppy_Disk_Android", "Shoe_Drugs", "Bird_Hnads", "Boat_Shower", "Allergies_Elevator", "Kitty_Website", "Crab_Body", "Dog_Puppy", "Leash_Toolbox", "Shower_Leash", "Laptop_Floppy_Disk", "Soda_Mail", "Settings_Book", "Male_Ice_cream_cone", "System_Elevator", "Websites_Websites", "Ice_cream_cone_Towel",
        "Shoe_Body"].map(s => s.slice(0, 10)).slice(0, USERS_COUNT)

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
            await joinRoom(userId, 'global-room')
            const imgUrl = await generateImageUrl(userId, bgColors[index % bgColors.length])
            for (let i = 0; i < CHEERS_PER_USER; i++) {
                console.log(`${i}/${CHEERS_PER_USER} ${userId}`)
                try {
                    await yum(userId, 'global-room', imgUrl)
                    await sleep(SLEEP_TIME);
                } catch (e) {
                    console.log(`${i}/${CHEERS_PER_USER} EXCEPTION ${userId} ${e}`)
                }
            }
        } catch (exception) {
            console.log(`${userId} ${exception}`)
        }
    }

    const all = userIds.map(joinAndSend)
    // const all = userIds.map((userId) => joinRoom(userId, 'global-room'))
    console.log("awaiting all jobs")
    await Promise.all(all)
    console.log("finished")
}

main()

if (TIMEOUT_SECONDS) {
    setTimeout(() => {
        process.exit()
    }, TIMEOUT_SECONDS * 1000)
}