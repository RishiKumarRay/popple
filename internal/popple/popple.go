package popple

import (
	"database/sql"
	"errors"

	"github.com/connorkuehl/popple/adapter"
)

type AddKarmaToEntitiesResult struct {
	Levels map[string]int64
	Err    error
}

type AddKarmaToEntitiesF chan AddKarmaToEntitiesResult

func AddKarmaToEntities(pl adapter.PersistenceLayer, serverID string, levels map[string]int64) AddKarmaToEntitiesF {
	f := make(chan AddKarmaToEntitiesResult, 1)
	go func() {
		updatedLevels := make(map[string]int64)
		for who, karma := range levels {
			if karma == 0 {
				continue
			}

			_ = pl.CreateEntity(serverID, who)
			updated, err := pl.AddKarmaToEntity(
				adapter.Entity{
					ServerID: serverID,
					Name:     who,
				},
				karma,
			)
			if err != nil {
				f <- AddKarmaToEntitiesResult{Err: err}
				return
			}

			updatedLevels[who] = int64(updated.Karma)
		}
		f <- AddKarmaToEntitiesResult{Levels: updatedLevels}
	}()
	return f
}

type GetConfigResult struct {
	C   adapter.Config
	Err error
}

type GetConfigF chan GetConfigResult

func GetConfig(pl adapter.PersistenceLayer, serverID string) GetConfigF {
	f := make(chan GetConfigResult, 1)
	go func() {
		c, e := pl.GetConfig(serverID)
		f <- GetConfigResult{
			C:   c,
			Err: e,
		}
	}()
	return f
}

type GetLevelsResult struct {
	Levels map[string]int64
	Err    error
}

type GetLevelsF chan GetLevelsResult

func GetLevels(pl adapter.PersistenceLayer, serverID string, bumps map[string]int64) GetLevelsF {
	f := make(chan GetLevelsResult, 1)
	go func() {
		levels := make(map[string]int64)
		for name := range bumps {
			entt, err := pl.GetEntity(serverID, name)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				f <- GetLevelsResult{Err: err}
				return
			}
			levels[name] = int64(entt.Karma)
		}
		f <- GetLevelsResult{Levels: levels}
	}()
	return f
}

type GetLeaderboardResult struct {
	Entries []adapter.LeaderboardEntry
	Err     error
}

type GetLeaderboardF chan GetLeaderboardResult

func GetLeaderboard(pl adapter.PersistenceLayer, serverID string, top bool, limit uint) GetLeaderboardF {
	f := make(chan GetLeaderboardResult, 1)
	go func() {
		board := func() ([]adapter.Entity, error) {
			if top {
				return pl.GetTopEntities(serverID, uint(limit))
			}
			return pl.GetBotEntities(serverID, uint(limit))
		}

		entities, err := board()
		if err != nil {
			f <- GetLeaderboardResult{Err: err}
			return
		}

		entries := make([]adapter.LeaderboardEntry, 0, len(entities))
		for _, entity := range entities {
			entries = append(entries, adapter.LeaderboardEntry{
				Name:  entity.Name,
				Karma: int64(entity.Karma),
			})
		}
		f <- GetLeaderboardResult{Entries: entries}
	}()
	return f
}

type SetAnnounceF chan error

func SetAnnounce(pl adapter.PersistenceLayer, serverID string, on bool) SetAnnounceF {
	f := make(chan error, 1)
	go func() {
		_ = pl.CreateConfig(serverID)
		f <- pl.PutConfig(adapter.Config{ServerID: serverID, NoAnnounce: !on})
	}()
	return f
}
