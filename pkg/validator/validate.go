package validator

/*
   Copyright (C) 2020 Daniel Gurney
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Validate validates a provided key.
func Validate(k string, v chan bool) {
	// Determine key type
	switch {
	case len(k) == 12 && k[4:5] == "-":
		ecd := elevencd{
			key: k,
		}
		ecd.validate(v)
	case len(k) == 11 && k[3:4] == "-":
		cd := cd{
			key: k,
		}
		cd.validate(v)
	case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
		oem := oem{
			key: k,
		}
		oem.validate(v)
	default:
		v <- false
	}
}
