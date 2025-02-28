package cache

import (
	"time"

	"codeberg.org/gruf/go-cache/v2"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

// AccountCache is a cache wrapper to provide URL and URI lookups for gtsmodel.Account
type AccountCache struct {
	cache cache.LookupCache[string, string, *gtsmodel.Account]
}

// NewAccountCache returns a new instantiated AccountCache object
func NewAccountCache() *AccountCache {
	c := &AccountCache{}
	c.cache = cache.NewLookup(cache.LookupCfg[string, string, *gtsmodel.Account]{
		RegisterLookups: func(lm *cache.LookupMap[string, string]) {
			lm.RegisterLookup("uri")
			lm.RegisterLookup("url")
		},

		AddLookups: func(lm *cache.LookupMap[string, string], acc *gtsmodel.Account) {
			if uri := acc.URI; uri != "" {
				lm.Set("uri", uri, acc.ID)
			}
			if url := acc.URL; url != "" {
				lm.Set("url", url, acc.ID)
			}
		},

		DeleteLookups: func(lm *cache.LookupMap[string, string], acc *gtsmodel.Account) {
			if uri := acc.URI; uri != "" {
				lm.Delete("uri", uri)
			}
			if url := acc.URL; url != "" {
				lm.Delete("url", url)
			}
		},
	})
	c.cache.SetTTL(time.Minute*5, false)
	c.cache.Start(time.Second * 10)
	return c
}

// GetByID attempts to fetch a account from the cache by its ID, you will receive a copy for thread-safety
func (c *AccountCache) GetByID(id string) (*gtsmodel.Account, bool) {
	return c.cache.Get(id)
}

// GetByURL attempts to fetch a account from the cache by its URL, you will receive a copy for thread-safety
func (c *AccountCache) GetByURL(url string) (*gtsmodel.Account, bool) {
	return c.cache.GetBy("url", url)
}

// GetByURI attempts to fetch a account from the cache by its URI, you will receive a copy for thread-safety
func (c *AccountCache) GetByURI(uri string) (*gtsmodel.Account, bool) {
	return c.cache.GetBy("uri", uri)
}

// Put places a account in the cache, ensuring that the object place is a copy for thread-safety
func (c *AccountCache) Put(account *gtsmodel.Account) {
	if account == nil || account.ID == "" {
		panic("invalid account")
	}
	c.cache.Set(account.ID, copyAccount(account))
}

// copyAccount performs a surface-level copy of account, only keeping attached IDs intact, not the objects.
// due to all the data being copied being 99% primitive types or strings (which are immutable and passed by ptr)
// this should be a relatively cheap process
func copyAccount(account *gtsmodel.Account) *gtsmodel.Account {
	return &gtsmodel.Account{
		ID:                      account.ID,
		Username:                account.Username,
		Domain:                  account.Domain,
		AvatarMediaAttachmentID: account.AvatarMediaAttachmentID,
		AvatarMediaAttachment:   nil,
		AvatarRemoteURL:         account.AvatarRemoteURL,
		HeaderMediaAttachmentID: account.HeaderMediaAttachmentID,
		HeaderMediaAttachment:   nil,
		HeaderRemoteURL:         account.HeaderRemoteURL,
		DisplayName:             account.DisplayName,
		Fields:                  account.Fields,
		Note:                    account.Note,
		NoteRaw:                 account.NoteRaw,
		Memorial:                account.Memorial,
		MovedToAccountID:        account.MovedToAccountID,
		CreatedAt:               account.CreatedAt,
		UpdatedAt:               account.UpdatedAt,
		Bot:                     account.Bot,
		Reason:                  account.Reason,
		Locked:                  account.Locked,
		Discoverable:            account.Discoverable,
		Privacy:                 account.Privacy,
		Sensitive:               account.Sensitive,
		Language:                account.Language,
		URI:                     account.URI,
		URL:                     account.URL,
		LastWebfingeredAt:       account.LastWebfingeredAt,
		InboxURI:                account.InboxURI,
		OutboxURI:               account.OutboxURI,
		FollowingURI:            account.FollowingURI,
		FollowersURI:            account.FollowersURI,
		FeaturedCollectionURI:   account.FeaturedCollectionURI,
		ActorType:               account.ActorType,
		AlsoKnownAs:             account.AlsoKnownAs,
		PrivateKey:              account.PrivateKey,
		PublicKey:               account.PublicKey,
		PublicKeyURI:            account.PublicKeyURI,
		SensitizedAt:            account.SensitizedAt,
		SilencedAt:              account.SilencedAt,
		SuspendedAt:             account.SuspendedAt,
		HideCollections:         account.HideCollections,
		SuspensionOrigin:        account.SuspensionOrigin,
	}
}
