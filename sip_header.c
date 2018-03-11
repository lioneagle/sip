/* C code produced by gperf version 3.0.1 */
/* Command-line: gperf -g -o -j1 -t -p -N GetSipHeaderIndex --ignore-case sip_headers.gperf  */
/* Computed positions: -k'1-2' */

#if !((' ' == 32) && ('!' == 33) && ('"' == 34) && ('#' == 35) \
      && ('%' == 37) && ('&' == 38) && ('\'' == 39) && ('(' == 40) \
      && (')' == 41) && ('*' == 42) && ('+' == 43) && (',' == 44) \
      && ('-' == 45) && ('.' == 46) && ('/' == 47) && ('0' == 48) \
      && ('1' == 49) && ('2' == 50) && ('3' == 51) && ('4' == 52) \
      && ('5' == 53) && ('6' == 54) && ('7' == 55) && ('8' == 56) \
      && ('9' == 57) && (':' == 58) && (';' == 59) && ('<' == 60) \
      && ('=' == 61) && ('>' == 62) && ('?' == 63) && ('A' == 65) \
      && ('B' == 66) && ('C' == 67) && ('D' == 68) && ('E' == 69) \
      && ('F' == 70) && ('G' == 71) && ('H' == 72) && ('I' == 73) \
      && ('J' == 74) && ('K' == 75) && ('L' == 76) && ('M' == 77) \
      && ('N' == 78) && ('O' == 79) && ('P' == 80) && ('Q' == 81) \
      && ('R' == 82) && ('S' == 83) && ('T' == 84) && ('U' == 85) \
      && ('V' == 86) && ('W' == 87) && ('X' == 88) && ('Y' == 89) \
      && ('Z' == 90) && ('[' == 91) && ('\\' == 92) && (']' == 93) \
      && ('^' == 94) && ('_' == 95) && ('a' == 97) && ('b' == 98) \
      && ('c' == 99) && ('d' == 100) && ('e' == 101) && ('f' == 102) \
      && ('g' == 103) && ('h' == 104) && ('i' == 105) && ('j' == 106) \
      && ('k' == 107) && ('l' == 108) && ('m' == 109) && ('n' == 110) \
      && ('o' == 111) && ('p' == 112) && ('q' == 113) && ('r' == 114) \
      && ('s' == 115) && ('t' == 116) && ('u' == 117) && ('v' == 118) \
      && ('w' == 119) && ('x' == 120) && ('y' == 121) && ('z' == 122) \
      && ('{' == 123) && ('|' == 124) && ('}' == 125) && ('~' == 126))
/* The character set is not based on ISO-646.  */
error "gperf generated tables don't work with this execution character set. Please report a bug to <bug-gnu-gperf@gnu.org>."
#endif

#line 1 "sip_headers.gperf"
struct resword { char *name; short header; };

#define TOTAL_KEYWORDS 44
#define MIN_WORD_LENGTH 1
#define MAX_WORD_LENGTH 19
#define MIN_HASH_VALUE 1
#define MAX_HASH_VALUE 44
/* maximum key range = 44, duplicates = 0 */

#ifndef GPERF_DOWNCASE
#define GPERF_DOWNCASE 1
static unsigned char gperf_downcase[256] =
  {
      0,   1,   2,   3,   4,   5,   6,   7,   8,   9,  10,  11,  12,  13,  14,
     15,  16,  17,  18,  19,  20,  21,  22,  23,  24,  25,  26,  27,  28,  29,
     30,  31,  32,  33,  34,  35,  36,  37,  38,  39,  40,  41,  42,  43,  44,
     45,  46,  47,  48,  49,  50,  51,  52,  53,  54,  55,  56,  57,  58,  59,
     60,  61,  62,  63,  64,  97,  98,  99, 100, 101, 102, 103, 104, 105, 106,
    107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121,
    122,  91,  92,  93,  94,  95,  96,  97,  98,  99, 100, 101, 102, 103, 104,
    105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119,
    120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134,
    135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149,
    150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164,
    165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179,
    180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194,
    195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209,
    210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224,
    225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239,
    240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254,
    255
  };
#endif

#ifndef GPERF_CASE_STRCMP
#define GPERF_CASE_STRCMP 1
static int
gperf_case_strcmp (s1, s2)
     register const char *s1;
     register const char *s2;
{
  for (;;)
    {
      unsigned char c1 = gperf_downcase[(unsigned char)*s1++];
      unsigned char c2 = gperf_downcase[(unsigned char)*s2++];
      if (c1 != 0 && c1 == c2)
        continue;
      return (int)c1 - (int)c2;
    }
}
#endif

#ifdef __GNUC__
__inline
#else
#ifdef __cplusplus
inline
#endif
#endif
static unsigned int
hash (str, len)
     register const char *str;
     register unsigned int len;
{
  static unsigned char asso_values[] =
    {
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45,  3, 43,  0, 33, 10,
      30, 45, 45,  4, 38, 37,  5, 27, 45,  2,
      45, 45,  1, 11, 28,  6, 26, 45, 36, 45,
      45, 45, 45, 45, 45, 45, 45,  3, 43,  0,
      33, 10, 30, 45, 45,  4, 38, 37,  5, 27,
      45,  2, 45, 45,  1, 11, 28,  6, 26, 45,
      36, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
      45, 45, 45, 45, 45, 45
    };
  register int hval = len;

  switch (hval)
    {
      default:
        hval += asso_values[(unsigned char)str[1]];
      /*FALLTHROUGH*/
      case 1:
        hval += asso_values[(unsigned char)str[0]];
        break;
    }
  return hval;
}

#ifdef __GNUC__
__inline
#endif
struct resword *
GetSipHeaderIndex (str, len)
     register const char *str;
     register unsigned int len;
{
  static struct resword wordlist[] =
    {
      {""},
#line 15 "sip_headers.gperf"
      {"c",                    ABNF_SIP_HDR_CONTENT_TYPE},
#line 35 "sip_headers.gperf"
      {"r",                    ABNF_SIP_HDR_REFER_TO},
#line 33 "sip_headers.gperf"
      {"o",                    ABNF_SIP_HDR_EVENT},
#line 37 "sip_headers.gperf"
      {"a",                    ABNF_SIP_HDR_ACCEPT_CONTACT},
#line 10 "sip_headers.gperf"
      {"i",                    ABNF_SIP_HDR_CALL_ID},
#line 13 "sip_headers.gperf"
      {"l",                    ABNF_SIP_HDR_CONTENT_LENGTH},
#line 31 "sip_headers.gperf"
      {"u",                    ABNF_SIP_HDR_ALLOW_EVENTS},
#line 19 "sip_headers.gperf"
      {"Route",                ABNF_SIP_HDR_ROUTE},
#line 16 "sip_headers.gperf"
      {"Contact",              ABNF_SIP_HDR_CONTACT},
#line 9 "sip_headers.gperf"
      {"Call-ID",              ABNF_SIP_HDR_CALL_ID},
#line 24 "sip_headers.gperf"
      {"e",                    ABNF_SIP_HDR_CONTENT_ENCODING},
#line 27 "sip_headers.gperf"
      {"s",                    ABNF_SIP_HDR_SUBJECT},
#line 22 "sip_headers.gperf"
      {"Allow",                ABNF_SIP_HDR_ALLOW},
#line 14 "sip_headers.gperf"
      {"Content-Type",         ABNF_SIP_HDR_CONTENT_TYPE},
#line 11 "sip_headers.gperf"
      {"CSeq",                 ABNF_SIP_HDR_CSEQ},
#line 12 "sip_headers.gperf"
      {"Content-Length",       ABNF_SIP_HDR_CONTENT_LENGTH},
#line 36 "sip_headers.gperf"
      {"Accept-Contact",       ABNF_SIP_HDR_ACCEPT_CONTACT},
#line 23 "sip_headers.gperf"
      {"Content-Encoding",     ABNF_SIP_HDR_CONTENT_ENCODING},
#line 34 "sip_headers.gperf"
      {"Refer-To",             ABNF_SIP_HDR_REFER_TO},
#line 30 "sip_headers.gperf"
      {"Allow-Events",         ABNF_SIP_HDR_ALLOW_EVENTS},
#line 21 "sip_headers.gperf"
      {"Content-Disposition",  ABNF_SIP_HDR_CONTENT_DISPOSITION},
#line 42 "sip_headers.gperf"
      {"Referred-By",          ABNF_SIP_HDR_REFERRED_BY},
#line 20 "sip_headers.gperf"
      {"Record-Route",         ABNF_SIP_HDR_RECORD_ROUTE},
#line 26 "sip_headers.gperf"
      {"Subject",              ABNF_SIP_HDR_SUBJECT},
#line 38 "sip_headers.gperf"
      {"Reject-Contact",       ABNF_SIP_HDR_REJECT_CONTACT},
#line 28 "sip_headers.gperf"
      {"Supported",            ABNF_SIP_HDR_SUPPORTED},
#line 8 "sip_headers.gperf"
      {"v",                    ABNF_SIP_HDR_VIA},
#line 17 "sip_headers.gperf"
      {"m",                    ABNF_SIP_HDR_CONTACT},
#line 6 "sip_headers.gperf"
      {"t",                    ABNF_SIP_HDR_TO},
#line 40 "sip_headers.gperf"
      {"Request-Disposition",  ABNF_SIP_HDR_REQUEST_DISPOSITION},
#line 4 "sip_headers.gperf"
      {"f",                    ABNF_SIP_HDR_FROM},
#line 5 "sip_headers.gperf"
      {"To",                   ABNF_SIP_HDR_TO},
#line 7 "sip_headers.gperf"
      {"Via",                  ABNF_SIP_HDR_VIA},
#line 41 "sip_headers.gperf"
      {"d",                    ABNF_SIP_HDR_REQUEST_DISPOSITION},
#line 3 "sip_headers.gperf"
      {"From",                 ABNF_SIP_HDR_FROM},
#line 44 "sip_headers.gperf"
      {"Session-Expires",      ABNF_SIP_HDR_SESSION_EXPIRES},
#line 45 "sip_headers.gperf"
      {"x",                    ABNF_SIP_HDR_SESSION_EXPIRES},
#line 29 "sip_headers.gperf"
      {"k",                    ABNF_SIP_HDR_SUPPORTED},
#line 39 "sip_headers.gperf"
      {"j",                    ABNF_SIP_HDR_REJECT_CONTACT},
#line 25 "sip_headers.gperf"
      {"Date",                 ABNF_SIP_HDR_DATE},
#line 32 "sip_headers.gperf"
      {"Event",                ABNF_SIP_HDR_EVENT},
#line 18 "sip_headers.gperf"
      {"Max-Forwards",         ABNF_SIP_HDR_MAX_FORWARDS},
#line 46 "sip_headers.gperf"
      {"MIME-Version",         ABNF_SIP_HDR_MIME_VERSION},
#line 43 "sip_headers.gperf"
      {"b",                    ABNF_SIP_HDR_REFERRED_BY}
    };

  if (len <= MAX_WORD_LENGTH && len >= MIN_WORD_LENGTH)
    {
      register int key = hash (str, len);

      if (key <= MAX_HASH_VALUE && key >= 0)
        {
          register const char *s = wordlist[key].name;

          if ((((unsigned char)*str ^ (unsigned char)*s) & ~32) == 0 && !gperf_case_strcmp (str, s))
            return &wordlist[key];
        }
    }
  return 0;
}
