package game

import (
	"embed"
	"encoding/json"
	"image/color"
	"io/fs"
	"log"
	"staticsheep/assets"
)

//go:embed levels/*.json
var levelFiles embed.FS

type levelObjectJSON struct {
	Type          string     `json:"type"`
	Name          string     `json:"name"`
	PositionPx    Vector     `json:"positionPx"`
	Radius        float64    `json:"radius"`
	WidthPx       float64    `json:"widthPx"`
	HeightPx      float64    `json:"heightPx"`
	Color         color.RGBA `json:"color"`
	BorderColor   color.RGBA `json:"borderColor"`
	BorderWidthPx float64    `json:"borderWidthPx"`
	Charge        float64    `json:"charge"`
	Sprite        string     `json:"sprite"`
}

func GetTotalNumLevels() int {
	numFiles := 0
	if err := fs.WalkDir(levelFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		numFiles++

		return nil
	}); err != nil {
		log.Println("Encountered error counting number of level files: ", err)
		numFiles = 0
	}

	return numFiles
}

// load all objects from JSON file.
// Using JSON so it's easy to edit
// would XML have been better for this?
func loadLevelFromJSON(path string, info *PlayerInfo) *Level {
	file, err := levelFiles.ReadFile(path)
	if err != nil {
		log.Printf("failed to read level file %s: %v", path, err)
		return nil
	}

	var objectData []levelObjectJSON
	if err := json.Unmarshal(file, &objectData); err != nil {
		log.Printf("failed to parse level file %s: %v", path, err)
		return nil
	}

	levelName := "NOLEVEL"
	backgroundColor := color.RGBA{10, 10, 20, 255}
	var player *Player
	var playerZone *PlayerZone
	var goal *GoalObject

	objects := make([]ChargedObject, 0, len(objectData))
	obstacles := make([]Obstacle, 0, len(objectData))
	for _, data := range objectData {
		// modify all positions to account for UI bar.
		// definitely better ways to have accomplished this, but it'll do for now
		data.PositionPx.Y += UiBarHeightPx

		// check type of data being read and create an object for it
		switch data.Type {
		case "obstacle":
			obstacles = append(obstacles, Obstacle{
				positionPx:    data.PositionPx,
				widthPx:       data.WidthPx,
				heightPx:      data.HeightPx,
				color:         data.Color,
				borderColor:   data.BorderColor,
				borderWidthPx: data.BorderWidthPx,
			})
		case "chargedObject":
			// load sprite
			sprite := assets.SheepRedSprite
			switch data.Sprite {
			case "sheep_blue":
				sprite = assets.SheepBlueSprite
			case "sheep_darkred":
				sprite = assets.SheepDarkRedSprite
			case "sheep_darkblue":
				sprite = assets.SheepDarkBlueSprite
			case "sheep_white":
				sprite = assets.SheepWhiteSprite
			case "sheep_red":
			default: // already initialized to red sheep sprite
			}
			objects = append(objects, ChargedObject{
				positionPx:    data.PositionPx,
				radius:        data.Radius,
				color:         data.Color,
				charge:        data.Charge,
				borderColor:   data.BorderColor,
				borderWidthPx: data.BorderWidthPx,
				sprite:        sprite,
			})
		case "goal":
			// always use white sheep for goal sprite (for now)
			goal = &GoalObject{
				ChargedObject: ChargedObject{
					positionPx:    data.PositionPx,
					radius:        data.Radius,
					color:         data.Color,
					charge:        data.Charge,
					borderColor:   data.BorderColor,
					borderWidthPx: data.BorderWidthPx,
					sprite:        assets.SheepWhiteSprite,
				},
			}
			// add just ChargedObject portion of goal to objects list
			objects = append(objects, goal.ChargedObject)
		case "playerZone":
			playerZone = &PlayerZone{
				positionPx:    data.PositionPx,
				widthPx:       data.WidthPx,
				heightPx:      data.HeightPx,
				color:         data.Color,
				borderColor:   data.BorderColor,
				borderWidthPx: data.BorderWidthPx,
			}
		case "player":
			player = NewPlayer(data.PositionPx)
		case "levelInfo":
			levelName = data.Name
			backgroundColor = data.Color
		default:

		}
	}

	// require a goal and player be defined
	if player == nil {
		log.Printf("Missing Player object definition in %s", path)
		return nil
	}
	if goal == nil {
		log.Printf("Missing Goal object definition in %s", path)
		return nil
	}

	return &Level{
		player:          player,
		goal:            goal,
		name:            levelName,
		objects:         objects,
		obstacles:       obstacles,
		playZone:        playerZone,
		backgroundColor: backgroundColor,
		info:            info,
	}
}
