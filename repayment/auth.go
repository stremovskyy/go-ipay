/*
 * Project: banker
 * File: auth.go (2/9/26, 10:12 AM)
 *
 * Copyright (C) Megakit Systems 2017-2026, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package repayment

// Auth holds the authentication details required for a Repayment API request.
type Auth struct {
	Login string `json:"login"`
	Time  string `json:"time"`
	Sign  string `json:"sign"`
}
