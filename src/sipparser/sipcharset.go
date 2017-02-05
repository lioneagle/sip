package sipparser

var charset0 = [256]uint32{
	0x0e000000, /* position 000 */
	0x0e000000, /* position 001 */
	0x0e000000, /* position 002 */
	0x0e000000, /* position 003 */
	0x0e000000, /* position 004 */
	0x0e000000, /* position 005 */
	0x0e000000, /* position 006 */
	0x0e000000, /* position 007 */
	0x0e000000, /* position 008 */
	0x0e001800, /* position 009 */
	0x0c001400, /* position 010 */
	0x0e800000, /* position 011 */
	0x0e000000, /* position 012 */
	0x0c001400, /* position 013 */
	0x0e000000, /* position 014 */
	0x0e000000, /* position 015 */
	0x0e000000, /* position 016 */
	0x0e000000, /* position 017 */
	0x0e000000, /* position 018 */
	0x0e000000, /* position 019 */
	0x0e000000, /* position 020 */
	0x0e000000, /* position 021 */
	0x0e000000, /* position 022 */
	0x0e000000, /* position 023 */
	0x0e000000, /* position 024 */
	0x0e000000, /* position 025 */
	0x0e000000, /* position 026 */
	0x0e000000, /* position 027 */
	0x0e000000, /* position 028 */
	0x0e000000, /* position 029 */
	0x0e000000, /* position 030 */
	0x0e000000, /* position 031 */
	0x0e801800, /* position 032  ' ', */
	0xff7ca000, /* position 033  '!', */
	0x0f802000, /* position 034  '"', */
	0x0e002000, /* position 035  '#', */
	0xfe3d2000, /* position 036  '$', */
	0x0f402000, /* position 037  '%', */
	0xfe3d2000, /* position 038  '&', */
	0xff7ca000, /* position 039  ''', */
	0xffbca000, /* position 040  '(', */
	0xffbca000, /* position 041  ')', */
	0xff7ca000, /* position 042  '*', */
	0xff7f2000, /* position 043  '+', */
	0x3ebd2000, /* position 044  ',', */
	0xff7ee000, /* position 045  '-', */
	0xff7ee000, /* position 046  '.', */
	0xde052000, /* position 047  '/', */
	0xff7ee391, /* position 048  '0', */
	0xff7ee391, /* position 049  '1', */
	0xff7ee391, /* position 050  '2', */
	0xff7ee391, /* position 051  '3', */
	0xff7ee391, /* position 052  '4', */
	0xff7ee391, /* position 053  '5', */
	0xff7ee391, /* position 054  '6', */
	0xff7ee391, /* position 055  '7', */
	0xff7ee391, /* position 056  '8', */
	0xff7ee391, /* position 057  '9', */
	0xcfbd2000, /* position 058  ':', */
	0x1ead2000, /* position 059  ';', */
	0x0f802000, /* position 060  '<', */
	0x3ebd2000, /* position 061  '=', */
	0x0f802000, /* position 062  '>', */
	0x1f8d2000, /* position 063  '?', */
	0x0ebd2000, /* position 064  '@', */
	0xff7ee35a, /* position 065  'A', */
	0xff7ee35a, /* position 066  'B', */
	0xff7ee35a, /* position 067  'C', */
	0xff7ee35a, /* position 068  'D', */
	0xff7ee35a, /* position 069  'E', */
	0xff7ee35a, /* position 070  'F', */
	0xff7ee01a, /* position 071  'G', */
	0xff7ee01a, /* position 072  'H', */
	0xff7ee01a, /* position 073  'I', */
	0xff7ee01a, /* position 074  'J', */
	0xff7ee01a, /* position 075  'K', */
	0xff7ee01a, /* position 076  'L', */
	0xff7ee01a, /* position 077  'M', */
	0xff7ee01a, /* position 078  'N', */
	0xff7ee01a, /* position 079  'O', */
	0xff7ee01a, /* position 080  'P', */
	0xff7ee01a, /* position 081  'Q', */
	0xff7ee01a, /* position 082  'R', */
	0xff7ee01a, /* position 083  'S', */
	0xff7ee01a, /* position 084  'T', */
	0xff7ee01a, /* position 085  'U', */
	0xff7ee01a, /* position 086  'V', */
	0xff7ee01a, /* position 087  'W', */
	0xff7ee01a, /* position 088  'X', */
	0xff7ee01a, /* position 089  'Y', */
	0xff7ee01a, /* position 090  'Z', */
	0xcf802000, /* position 091  '[', */
	0x0f802000, /* position 092  '\', */
	0xcf802000, /* position 093  ']', */
	0x0e002000, /* position 094  '^', */
	0xff7ca000, /* position 095  '_', */
	0x0f402000, /* position 096  '`', */
	0xff7ee2b6, /* position 097  'a', */
	0xff7ee2b6, /* position 098  'b', */
	0xff7ee2b6, /* position 099  'c', */
	0xff7ee2b6, /* position 100  'd', */
	0xff7ee2b6, /* position 101  'e', */
	0xff7ee2b6, /* position 102  'f', */
	0xff7ee016, /* position 103  'g', */
	0xff7ee016, /* position 104  'h', */
	0xff7ee016, /* position 105  'i', */
	0xff7ee016, /* position 106  'j', */
	0xff7ee016, /* position 107  'k', */
	0xff7ee016, /* position 108  'l', */
	0xff7ee016, /* position 109  'm', */
	0xff7ee016, /* position 110  'n', */
	0xff7ee016, /* position 111  'o', */
	0xff7ee016, /* position 112  'p', */
	0xff7ee016, /* position 113  'q', */
	0xff7ee016, /* position 114  'r', */
	0xff7ee016, /* position 115  's', */
	0xff7ee016, /* position 116  't', */
	0xff7ee016, /* position 117  'u', */
	0xff7ee016, /* position 118  'v', */
	0xff7ee016, /* position 119  'w', */
	0xff7ee016, /* position 120  'x', */
	0xff7ee016, /* position 121  'y', */
	0xff7ee016, /* position 122  'z', */
	0x0f802000, /* position 123  '{', */
	0x0e002000, /* position 124  '|', */
	0x0f802000, /* position 125  '}', */
	0xff7ca000, /* position 126  '~', */
	0x0e000000, /* position 127 */
	0x00002000, /* position 128 */
	0x00002000, /* position 129 */
	0x00002000, /* position 130 */
	0x00002000, /* position 131 */
	0x00002000, /* position 132 */
	0x00002000, /* position 133 */
	0x00002000, /* position 134 */
	0x00002000, /* position 135 */
	0x00002000, /* position 136 */
	0x00002000, /* position 137 */
	0x00002000, /* position 138 */
	0x00002000, /* position 139 */
	0x00002000, /* position 140 */
	0x00002000, /* position 141 */
	0x00002000, /* position 142 */
	0x00002000, /* position 143 */
	0x00002000, /* position 144 */
	0x00002000, /* position 145 */
	0x00002000, /* position 146 */
	0x00002000, /* position 147 */
	0x00002000, /* position 148 */
	0x00002000, /* position 149 */
	0x00002000, /* position 150 */
	0x00002000, /* position 151 */
	0x00002000, /* position 152 */
	0x00002000, /* position 153 */
	0x00002000, /* position 154 */
	0x00002000, /* position 155 */
	0x00002000, /* position 156 */
	0x00002000, /* position 157 */
	0x00002000, /* position 158 */
	0x00002000, /* position 159 */
	0x00002000, /* position 160 */
	0x00002000, /* position 161 */
	0x00002000, /* position 162 */
	0x00002000, /* position 163 */
	0x00002000, /* position 164 */
	0x00002000, /* position 165 */
	0x00002000, /* position 166 */
	0x00002000, /* position 167 */
	0x00002000, /* position 168 */
	0x00002000, /* position 169 */
	0x00002000, /* position 170 */
	0x00002000, /* position 171 */
	0x00002000, /* position 172 */
	0x00002000, /* position 173 */
	0x00002000, /* position 174 */
	0x00002000, /* position 175 */
	0x00002000, /* position 176 */
	0x00002000, /* position 177 */
	0x00002000, /* position 178 */
	0x00002000, /* position 179 */
	0x00002000, /* position 180 */
	0x00002000, /* position 181 */
	0x00002000, /* position 182 */
	0x00002000, /* position 183 */
	0x00002000, /* position 184 */
	0x00002000, /* position 185 */
	0x00002000, /* position 186 */
	0x00002000, /* position 187 */
	0x00002000, /* position 188 */
	0x00002000, /* position 189 */
	0x00002000, /* position 190 */
	0x00002000, /* position 191 */
	0x0c002000, /* position 192 */
	0x0c002000, /* position 193 */
	0x0c002000, /* position 194 */
	0x0c002000, /* position 195 */
	0x0c002000, /* position 196 */
	0x0c002000, /* position 197 */
	0x0c002000, /* position 198 */
	0x0c002000, /* position 199 */
	0x0c002000, /* position 200 */
	0x0c002000, /* position 201 */
	0x0c002000, /* position 202 */
	0x0c002000, /* position 203 */
	0x0c002000, /* position 204 */
	0x0c002000, /* position 205 */
	0x0c002000, /* position 206 */
	0x0c002000, /* position 207 */
	0x0c002000, /* position 208 */
	0x0c002000, /* position 209 */
	0x0c002000, /* position 210 */
	0x0c002000, /* position 211 */
	0x0c002000, /* position 212 */
	0x0c002000, /* position 213 */
	0x0c002000, /* position 214 */
	0x0c002000, /* position 215 */
	0x0c002000, /* position 216 */
	0x0c002000, /* position 217 */
	0x0c002000, /* position 218 */
	0x0c002000, /* position 219 */
	0x0c002000, /* position 220 */
	0x0c002000, /* position 221 */
	0x0c002000, /* position 222 */
	0x0c002000, /* position 223 */
	0x0c002000, /* position 224 */
	0x0c002000, /* position 225 */
	0x0c002000, /* position 226 */
	0x0c002000, /* position 227 */
	0x0c002000, /* position 228 */
	0x0c002000, /* position 229 */
	0x0c002000, /* position 230 */
	0x0c002000, /* position 231 */
	0x0c002000, /* position 232 */
	0x0c002000, /* position 233 */
	0x0c002000, /* position 234 */
	0x0c002000, /* position 235 */
	0x0c002000, /* position 236 */
	0x0c002000, /* position 237 */
	0x0c002000, /* position 238 */
	0x0c002000, /* position 239 */
	0x0c002000, /* position 240 */
	0x0c002000, /* position 241 */
	0x0c002000, /* position 242 */
	0x0c002000, /* position 243 */
	0x0c002000, /* position 244 */
	0x0c002000, /* position 245 */
	0x0c002000, /* position 246 */
	0x0c002000, /* position 247 */
	0x0c002000, /* position 248 */
	0x0c002000, /* position 249 */
	0x0c002000, /* position 250 */
	0x0c002000, /* position 251 */
	0x0c002000, /* position 252 */
	0x0c002000, /* position 253 */
	0x00000000, /* position 254 */
	0x00000000, /* position 255 */
}
var charset1 = [256]uint32{
	0x00000000, /* position 000 */
	0x00000000, /* position 001 */
	0x00000000, /* position 002 */
	0x00000000, /* position 003 */
	0x00000000, /* position 004 */
	0x00000000, /* position 005 */
	0x00000000, /* position 006 */
	0x00000000, /* position 007 */
	0x00000000, /* position 008 */
	0x00000000, /* position 009 */
	0x00000000, /* position 010 */
	0x00000004, /* position 011 */
	0x00000000, /* position 012 */
	0x00000000, /* position 013 */
	0x00000000, /* position 014 */
	0x00000000, /* position 015 */
	0x00000000, /* position 016 */
	0x00000000, /* position 017 */
	0x00000000, /* position 018 */
	0x00000000, /* position 019 */
	0x00000000, /* position 020 */
	0x00000000, /* position 021 */
	0x00000000, /* position 022 */
	0x00000000, /* position 023 */
	0x00000000, /* position 024 */
	0x00000000, /* position 025 */
	0x00000000, /* position 026 */
	0x00000000, /* position 027 */
	0x00000000, /* position 028 */
	0x00000000, /* position 029 */
	0x00000000, /* position 030 */
	0x00000000, /* position 031 */
	0x00000004, /* position 032  ' ', */
	0x00000047, /* position 033  '!', */
	0x00000000, /* position 034  '"', */
	0x00000000, /* position 035  '#', */
	0x00000047, /* position 036  '$', */
	0x00000000, /* position 037  '%', */
	0x00000044, /* position 038  '&', */
	0x00000047, /* position 039  ''', */
	0x000000df, /* position 040  '(', */
	0x000000df, /* position 041  ')', */
	0x00000047, /* position 042  '*', */
	0x00000047, /* position 043  '+', */
	0x00000004, /* position 044  ',', */
	0x000000ff, /* position 045  '-', */
	0x000000df, /* position 046  '.', */
	0x00000047, /* position 047  '/', */
	0x0000007f, /* position 048  '0', */
	0x0000007f, /* position 049  '1', */
	0x0000007f, /* position 050  '2', */
	0x0000007f, /* position 051  '3', */
	0x0000007f, /* position 052  '4', */
	0x0000007f, /* position 053  '5', */
	0x0000007f, /* position 054  '6', */
	0x0000007f, /* position 055  '7', */
	0x0000007f, /* position 056  '8', */
	0x0000007f, /* position 057  '9', */
	0x00000047, /* position 058  ':', */
	0x00000004, /* position 059  ';', */
	0x00000000, /* position 060  '<', */
	0x00000004, /* position 061  '=', */
	0x00000000, /* position 062  '>', */
	0x00000007, /* position 063  '?', */
	0x00000004, /* position 064  '@', */
	0x00000077, /* position 065  'A', */
	0x00000077, /* position 066  'B', */
	0x00000077, /* position 067  'C', */
	0x00000077, /* position 068  'D', */
	0x00000077, /* position 069  'E', */
	0x00000077, /* position 070  'F', */
	0x00000067, /* position 071  'G', */
	0x00000067, /* position 072  'H', */
	0x00000067, /* position 073  'I', */
	0x00000067, /* position 074  'J', */
	0x00000067, /* position 075  'K', */
	0x00000067, /* position 076  'L', */
	0x00000067, /* position 077  'M', */
	0x00000067, /* position 078  'N', */
	0x00000067, /* position 079  'O', */
	0x00000067, /* position 080  'P', */
	0x00000067, /* position 081  'Q', */
	0x00000067, /* position 082  'R', */
	0x00000067, /* position 083  'S', */
	0x00000067, /* position 084  'T', */
	0x00000067, /* position 085  'U', */
	0x00000067, /* position 086  'V', */
	0x00000067, /* position 087  'W', */
	0x00000067, /* position 088  'X', */
	0x00000067, /* position 089  'Y', */
	0x00000067, /* position 090  'Z', */
	0x00000043, /* position 091  '[', */
	0x00000000, /* position 092  '\', */
	0x00000043, /* position 093  ']', */
	0x00000000, /* position 094  '^', */
	0x00000047, /* position 095  '_', */
	0x00000000, /* position 096  '`', */
	0x00000077, /* position 097  'a', */
	0x00000077, /* position 098  'b', */
	0x00000077, /* position 099  'c', */
	0x00000077, /* position 100  'd', */
	0x00000077, /* position 101  'e', */
	0x00000077, /* position 102  'f', */
	0x00000067, /* position 103  'g', */
	0x00000067, /* position 104  'h', */
	0x00000067, /* position 105  'i', */
	0x00000067, /* position 106  'j', */
	0x00000067, /* position 107  'k', */
	0x00000067, /* position 108  'l', */
	0x00000067, /* position 109  'm', */
	0x00000067, /* position 110  'n', */
	0x00000067, /* position 111  'o', */
	0x00000067, /* position 112  'p', */
	0x00000067, /* position 113  'q', */
	0x00000067, /* position 114  'r', */
	0x00000067, /* position 115  's', */
	0x00000067, /* position 116  't', */
	0x00000067, /* position 117  'u', */
	0x00000067, /* position 118  'v', */
	0x00000067, /* position 119  'w', */
	0x00000067, /* position 120  'x', */
	0x00000067, /* position 121  'y', */
	0x00000067, /* position 122  'z', */
	0x00000000, /* position 123  '{', */
	0x00000000, /* position 124  '|', */
	0x00000000, /* position 125  '}', */
	0x00000047, /* position 126  '~', */
	0x00000000, /* position 127 */
	0x00000004, /* position 128 */
	0x00000004, /* position 129 */
	0x00000004, /* position 130 */
	0x00000004, /* position 131 */
	0x00000004, /* position 132 */
	0x00000004, /* position 133 */
	0x00000004, /* position 134 */
	0x00000004, /* position 135 */
	0x00000004, /* position 136 */
	0x00000004, /* position 137 */
	0x00000004, /* position 138 */
	0x00000004, /* position 139 */
	0x00000004, /* position 140 */
	0x00000004, /* position 141 */
	0x00000004, /* position 142 */
	0x00000004, /* position 143 */
	0x00000004, /* position 144 */
	0x00000004, /* position 145 */
	0x00000004, /* position 146 */
	0x00000004, /* position 147 */
	0x00000004, /* position 148 */
	0x00000004, /* position 149 */
	0x00000004, /* position 150 */
	0x00000004, /* position 151 */
	0x00000004, /* position 152 */
	0x00000004, /* position 153 */
	0x00000004, /* position 154 */
	0x00000004, /* position 155 */
	0x00000004, /* position 156 */
	0x00000004, /* position 157 */
	0x00000004, /* position 158 */
	0x00000004, /* position 159 */
	0x00000004, /* position 160 */
	0x00000004, /* position 161 */
	0x00000004, /* position 162 */
	0x00000004, /* position 163 */
	0x00000004, /* position 164 */
	0x00000004, /* position 165 */
	0x00000004, /* position 166 */
	0x00000004, /* position 167 */
	0x00000004, /* position 168 */
	0x00000004, /* position 169 */
	0x00000004, /* position 170 */
	0x00000004, /* position 171 */
	0x00000004, /* position 172 */
	0x00000004, /* position 173 */
	0x00000004, /* position 174 */
	0x00000004, /* position 175 */
	0x00000004, /* position 176 */
	0x00000004, /* position 177 */
	0x00000004, /* position 178 */
	0x00000004, /* position 179 */
	0x00000004, /* position 180 */
	0x00000004, /* position 181 */
	0x00000004, /* position 182 */
	0x00000004, /* position 183 */
	0x00000004, /* position 184 */
	0x00000004, /* position 185 */
	0x00000004, /* position 186 */
	0x00000004, /* position 187 */
	0x00000004, /* position 188 */
	0x00000004, /* position 189 */
	0x00000004, /* position 190 */
	0x00000004, /* position 191 */
	0x00000004, /* position 192 */
	0x00000004, /* position 193 */
	0x00000004, /* position 194 */
	0x00000004, /* position 195 */
	0x00000004, /* position 196 */
	0x00000004, /* position 197 */
	0x00000004, /* position 198 */
	0x00000004, /* position 199 */
	0x00000004, /* position 200 */
	0x00000004, /* position 201 */
	0x00000004, /* position 202 */
	0x00000004, /* position 203 */
	0x00000004, /* position 204 */
	0x00000004, /* position 205 */
	0x00000004, /* position 206 */
	0x00000004, /* position 207 */
	0x00000004, /* position 208 */
	0x00000004, /* position 209 */
	0x00000004, /* position 210 */
	0x00000004, /* position 211 */
	0x00000004, /* position 212 */
	0x00000004, /* position 213 */
	0x00000004, /* position 214 */
	0x00000004, /* position 215 */
	0x00000004, /* position 216 */
	0x00000004, /* position 217 */
	0x00000004, /* position 218 */
	0x00000004, /* position 219 */
	0x00000004, /* position 220 */
	0x00000004, /* position 221 */
	0x00000004, /* position 222 */
	0x00000004, /* position 223 */
	0x00000004, /* position 224 */
	0x00000004, /* position 225 */
	0x00000004, /* position 226 */
	0x00000004, /* position 227 */
	0x00000004, /* position 228 */
	0x00000004, /* position 229 */
	0x00000004, /* position 230 */
	0x00000004, /* position 231 */
	0x00000004, /* position 232 */
	0x00000004, /* position 233 */
	0x00000004, /* position 234 */
	0x00000004, /* position 235 */
	0x00000004, /* position 236 */
	0x00000004, /* position 237 */
	0x00000004, /* position 238 */
	0x00000004, /* position 239 */
	0x00000004, /* position 240 */
	0x00000004, /* position 241 */
	0x00000004, /* position 242 */
	0x00000004, /* position 243 */
	0x00000004, /* position 244 */
	0x00000004, /* position 245 */
	0x00000004, /* position 246 */
	0x00000004, /* position 247 */
	0x00000004, /* position 248 */
	0x00000004, /* position 249 */
	0x00000004, /* position 250 */
	0x00000004, /* position 251 */
	0x00000004, /* position 252 */
	0x00000004, /* position 253 */
	0x00000000, /* position 254 */
	0x00000000, /* position 255 */
}

const MASK_DIGIT uint32 = (0x00000001)
const MASK_ALPHA uint32 = (0x00000002)
const MASK_LOWER uint32 = (0x00000004)
const MASK_UPPER uint32 = (0x00000008)
const MASK_ALPHANUM uint32 = (0x00000010)
const MASK_LOWER_HEX_ALPHA uint32 = (0x00000020)
const MASK_UPPER_HEX_ALPHA uint32 = (0x00000040)
const MASK_LOWER_HEX uint32 = (0x00000080)
const MASK_UPPER_HEX uint32 = (0x00000100)
const MASK_HEX uint32 = (0x00000200)
const MASK_CRLF_CHAR uint32 = (0x00000400)
const MASK_WSP_CHAR uint32 = (0x00000800)
const MASK_LWS_CHAR uint32 = (0x00001000)
const MASK_UTF8_CHAR uint32 = (0x00002000)
const MASK_HOSTNAME uint32 = (0x00004000)
const MASK_URI_UNRESERVED uint32 = (0x00008000)
const MASK_URI_RESERVED uint32 = (0x00010000)
const MASK_URI_SCHEME uint32 = (0x00020000)
const MASK_URI_URIC uint32 = (0x00040000)
const MASK_URI_URIC_NO_SLASH uint32 = (0x00080000)
const MASK_URI_PCHAR uint32 = (0x00100000)
const MASK_URI_REG_NAME uint32 = (0x00200000)
const MASK_SIP_TOKEN uint32 = (0x00400000)
const MASK_SIP_SEPARATORS uint32 = (0x00800000)
const MASK_SIP_WORD uint32 = (0x01000000)
const MASK_SIP_QUOTED_PAIR uint32 = (0x02000000)
const MASK_SIP_QUOTED_STRING uint32 = (0x04000000)
const MASK_SIP_COMMENT uint32 = (0x08000000)
const MASK_SIP_USER uint32 = (0x10000000)
const MASK_SIP_PASSWORD uint32 = (0x20000000)
const MASK_SIP_PNAME uint32 = (0x40000000)
const MASK_SIP_PVALUE uint32 = (0x80000000)
const MASK_SIP_HNAME uint32 = (0x00000001)
const MASK_SIP_HVALUE uint32 = (0x00000002)
const MASK_SIP_REASON_PHRASE uint32 = (0x00000004)
const MASK_TEL_PHONE_DIGIT uint32 = (0x00000008)
const MASK_TEL_PHONE_DIGIT_HEX uint32 = (0x00000010)
const MASK_TEL_PNAME uint32 = (0x00000020)
const MASK_TEL_PVALUE uint32 = (0x00000040)
const MASK_TEL_VISUAL_SPERATOR uint32 = (0x00000080)

func IsDigit(ch byte) bool             { return (charset0[ch] & MASK_DIGIT) != 0 }
func IsAlpha(ch byte) bool             { return (charset0[ch] & MASK_ALPHA) != 0 }
func IsLower(ch byte) bool             { return (charset0[ch] & MASK_LOWER) != 0 }
func IsUpper(ch byte) bool             { return (charset0[ch] & MASK_UPPER) != 0 }
func IsAlphanum(ch byte) bool          { return (charset0[ch] & MASK_ALPHANUM) != 0 }
func IsLowerHexAlpha(ch byte) bool     { return (charset0[ch] & MASK_LOWER_HEX_ALPHA) != 0 }
func IsUpperHexAlpha(ch byte) bool     { return (charset0[ch] & MASK_UPPER_HEX_ALPHA) != 0 }
func IsLowerHex(ch byte) bool          { return (charset0[ch] & MASK_LOWER_HEX) != 0 }
func IsUpperHex(ch byte) bool          { return (charset0[ch] & MASK_UPPER_HEX) != 0 }
func IsHex(ch byte) bool               { return (charset0[ch] & MASK_HEX) != 0 }
func IsCrlfChar(ch byte) bool          { return (charset0[ch] & MASK_CRLF_CHAR) != 0 }
func IsWspChar(ch byte) bool           { return (charset0[ch] & MASK_WSP_CHAR) != 0 }
func IsLwsChar(ch byte) bool           { return (charset0[ch] & MASK_LWS_CHAR) != 0 }
func IsUtf8Char(ch byte) bool          { return (charset0[ch] & MASK_UTF8_CHAR) != 0 }
func IsHostname(ch byte) bool          { return (charset0[ch] & MASK_HOSTNAME) != 0 }
func IsUriUnreserved(ch byte) bool     { return (charset0[ch] & MASK_URI_UNRESERVED) != 0 }
func IsUriReserved(ch byte) bool       { return (charset0[ch] & MASK_URI_RESERVED) != 0 }
func IsUriScheme(ch byte) bool         { return (charset0[ch] & MASK_URI_SCHEME) != 0 }
func IsUriUric(ch byte) bool           { return (charset0[ch] & MASK_URI_URIC) != 0 }
func IsUriUricNoSlash(ch byte) bool    { return (charset0[ch] & MASK_URI_URIC_NO_SLASH) != 0 }
func IsUriPchar(ch byte) bool          { return (charset0[ch] & MASK_URI_PCHAR) != 0 }
func IsUriRegName(ch byte) bool        { return (charset0[ch] & MASK_URI_REG_NAME) != 0 }
func IsSipToken(ch byte) bool          { return (charset0[ch] & MASK_SIP_TOKEN) != 0 }
func IsSipSeparators(ch byte) bool     { return (charset0[ch] & MASK_SIP_SEPARATORS) != 0 }
func IsSipWord(ch byte) bool           { return (charset0[ch] & MASK_SIP_WORD) != 0 }
func IsSipQuotedPair(ch byte) bool     { return (charset0[ch] & MASK_SIP_QUOTED_PAIR) != 0 }
func IsSipQuotedString(ch byte) bool   { return (charset0[ch] & MASK_SIP_QUOTED_STRING) != 0 }
func IsSipComment(ch byte) bool        { return (charset0[ch] & MASK_SIP_COMMENT) != 0 }
func IsSipUser(ch byte) bool           { return (charset0[ch] & MASK_SIP_USER) != 0 }
func IsSipPassword(ch byte) bool       { return (charset0[ch] & MASK_SIP_PASSWORD) != 0 }
func IsSipPname(ch byte) bool          { return (charset0[ch] & MASK_SIP_PNAME) != 0 }
func IsSipPvalue(ch byte) bool         { return (charset0[ch] & MASK_SIP_PVALUE) != 0 }
func IsSipHname(ch byte) bool          { return (charset1[ch] & MASK_SIP_HNAME) != 0 }
func IsSipHvalue(ch byte) bool         { return (charset1[ch] & MASK_SIP_HVALUE) != 0 }
func IsSipReasonPhrase(ch byte) bool   { return (charset1[ch] & MASK_SIP_REASON_PHRASE) != 0 }
func IsTelPhoneDigit(ch byte) bool     { return (charset1[ch] & MASK_TEL_PHONE_DIGIT) != 0 }
func IsTelPhoneDigitHex(ch byte) bool  { return (charset1[ch] & MASK_TEL_PHONE_DIGIT_HEX) != 0 }
func IsTelPname(ch byte) bool          { return (charset1[ch] & MASK_TEL_PNAME) != 0 }
func IsTelPvalue(ch byte) bool         { return (charset1[ch] & MASK_TEL_PVALUE) != 0 }
func IsTelVisualSperator(ch byte) bool { return (charset1[ch] & MASK_TEL_VISUAL_SPERATOR) != 0 }
