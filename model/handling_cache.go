package model

//CacheObj - CacheObj
type CacheObj struct {
	Type        string
	ID          string
	Additional  string
	Additional1 string
}

//CacheChn - channel
var CacheChn = make(chan CacheObj)

// HandleCaching - Handle all cache changes and sync process
func HandleCaching() {
	for {
		select {
		case message := <-CacheChn:
			mtype := message.Type
			switch mtype {
			case "hotel":
				if UpdateHotelOnList(message.ID) {
					UpdateHotelTag(message.ID)
				}
				break
			case "tax":
				UpdateHotelTax(message.ID)
				break
			case "roomType":
				UpdateAllRoomType(message.ID, message.Additional)
				break
			case "roomTypeAmenity":
				UpdateRoomAmenity(message.ID, message.Additional)
				break
			case "roomTypeImages":
				UpdateRoomImage(message.ID, message.Additional)
				break
			case "ratePlanDetails":
				AddUpdateRateplanDetails(message.ID, message.Additional)
				break
			case "updateDeals":
				UpdateRatePlanDeals(message.ID, message.Additional, message.Additional1)
			//Run RPC call of change property status in background
			case "changePropertyStatus":
				UpdatePropertyFlag(message.ID)
				break
			case "updateHotelWithProperty":
				UpdateHotelWithProperty(message.ID)
				break
			default:
				break
			}
		}
	}
}
