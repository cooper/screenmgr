/*
 *  vncauth.c
 *
 *  A consolidation of the following files, with modifications:
 *       d3des.h
 *       d3des.c
 *       vncauth.h
 *       vncauth.c
 *       This implementation of D3DES is released to the public domain.
 *   The code from vncauth/vncpasswd is licensed under GPL.
 */

 /*
  * This is D3DES (V5.09) by Richard Outerbridge with the double and
  * triple-length support removed for use in VNC.  Also the bytebit[] array
  * has been reversed so that the most significant bit in each byte of the
  * key is ignored, not the least significant.
  *
  * These changes are:
  * Copyright (C) 1999 AT&T Laboratories Cambridge. All Rights Reserved.
  *
  * This software is distributed in the hope that it will be useful,
  * but WITHOUT ANY WARRANTY; without even the implied warranty of
  * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
  */

 /* D3DES (V5.09) -
  *
  * A portable, public domain, version of the Data Encryption Standard.
  *
  * Written with Symantec's THINK (Lightspeed) C by Richard Outerbridge.
  * Thanks to: Dan Hoey for his excellent Initial and Inverse permutation
  * code;  Jim Gillogly & Phil Karn for the DES key schedule code; Dennis
  * Ferguson, Eric Young and Dana How for comparing notes; and Ray Lau,
  * for humouring me on.
  *
  * Copyright (c) 1988,1989,1990,1991,1992 by Richard Outerbridge.
  * (GEnie : OUTER; CIS : [71755,204]) Graven Imagery, 1992.
  */

 /*
  * This is D3DES (V5.09) by Richard Outerbridge with the double and
  * triple-length support removed for use in VNC.
  *
  * These changes are
  * Copyright (C) 1999 AT&T Laboratories Cambridge. All Rights Reserved.
  *
  * This software is distributed in the hope that it will be useful,
  * but WITHOUT ANY WARRANTY; without even the implied warranty of
  * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
  */

 /* d3des.h -
  *
  *      Headers and defines for d3des.c
  *      Graven Imagery, 1992.
  *
  * Copyright (c) 1988,1989,1990,1991,1992 by Richard Outerbridge
  *      (GEnie : OUTER; CIS : [71755,204])
  */

 #define EN0     0       /* MODE == encrypt */
 #define DE1     1       /* MODE == decrypt */

 /* d3des.h V5.09 rwo 9208.04 15:06 Graven Imagery
  ********************************************************************/

 static void scrunch(unsigned char *, unsigned long *);
 static void unscrun(unsigned long *, unsigned char *);
 static void desfunc(unsigned long *, unsigned long *);
 static void cookey(unsigned long *);

 static unsigned long KnL[32] = { 0L };
 /*static unsigned long KnR[32] = { 0L };*/
 /*static unsigned long Kn3[32] = { 0L };*/
 /*static unsigned char Df_Key[24] = {
  0x01,0x23,0x45,0x67,0x89,0xab,0xcd,0xef,
  0xfe,0xdc,0xba,0x98,0x76,0x54,0x32,0x10,
  0x89,0xab,0xcd,0xef,0x01,0x23,0x45,0x67 };*/

 static unsigned short bytebit[8]        = {
     01, 02, 04, 010, 020, 040, 0100, 0200 };

 static unsigned long bigbyte[24] = {
     0x800000L,      0x400000L,      0x200000L,      0x100000L,
     0x80000L,       0x40000L,       0x20000L,       0x10000L,
     0x8000L,        0x4000L,        0x2000L,        0x1000L,
     0x800L,         0x400L,         0x200L,         0x100L,
     0x80L,          0x40L,          0x20L,          0x10L,
     0x8L,           0x4L,           0x2L,           0x1L    };

 /* Use the key schedule specified in the Standard (ANSI X3.92-1981). */

 static unsigned char pc1[56] = {
     56, 48, 40, 32, 24, 16,  8,      0, 57, 49, 41, 33, 25, 17,
     9,  1, 58, 50, 42, 34, 26,     18, 10,  2, 59, 51, 43, 35,
     62, 54, 46, 38, 30, 22, 14,      6, 61, 53, 45, 37, 29, 21,
     13,  5, 60, 52, 44, 36, 28,     20, 12,  4, 27, 19, 11,  3 };

 static unsigned char totrot[16] = {
     1,2,4,6,8,10,12,14,15,17,19,21,23,25,27,28 };

 static unsigned char pc2[48] = {
     13, 16, 10, 23,  0,  4,  2, 27, 14,  5, 20,  9,
     22, 18, 11,  3, 25,  7, 15,  6, 26, 19, 12,  1,
     40, 51, 30, 36, 46, 54, 29, 39, 50, 44, 32, 47,
     43, 48, 38, 55, 33, 52, 45, 41, 49, 35, 28, 31 };

 void deskey(key, edf)   /* Thanks to James Gillogly & Phil Karn! */
 unsigned char *key;
 int edf;
 {
     register int i, j, l, m, n;
     unsigned char pc1m[56], pcr[56];
     unsigned long kn[32];

     for ( j = 0; j < 56; j++ ) {
         l = pc1[j];
         m = l & 07;
         pc1m[j] = (key[l >> 3] & bytebit[m]) ? 1 : 0;
     }
     for( i = 0; i < 16; i++ ) {
         if( edf == DE1 ) m = (15 - i) << 1;
         else m = i << 1;
         n = m + 1;
         kn[m] = kn[n] = 0L;
         for( j = 0; j < 28; j++ ) {
             l = j + totrot[i];
             if( l < 28 ) pcr[j] = pc1m[l];
             else pcr[j] = pc1m[l - 28];
         }
         for( j = 28; j < 56; j++ ) {
             l = j + totrot[i];
             if( l < 56 ) pcr[j] = pc1m[l];
             else pcr[j] = pc1m[l - 28];
         }
         for( j = 0; j < 24; j++ ) {
             if( pcr[pc2[j]] ) kn[m] |= bigbyte[j];
             if( pcr[pc2[j+24]] ) kn[n] |= bigbyte[j];
         }
     }
     cookey(kn);
     return;
 }

 void usekey(from)
 register unsigned long *from;
 {
     register unsigned long *to, *endp;

     to = KnL, endp = &KnL[32];
     while( to < endp ) *to++ = *from++;
         return;
 }

 static void cookey(raw1)
 register unsigned long *raw1;
 {
     register unsigned long *cook, *raw0;
     unsigned long dough[32];
     register int i;

     cook = dough;
     for( i = 0; i < 16; i++, raw1++ ) {
         raw0 = raw1++;
         *cook    = (*raw0 & 0x00fc0000L) << 6;
         *cook   |= (*raw0 & 0x00000fc0L) << 10;
         *cook   |= (*raw1 & 0x00fc0000L) >> 10;
         *cook++ |= (*raw1 & 0x00000fc0L) >> 6;
         *cook    = (*raw0 & 0x0003f000L) << 12;
         *cook   |= (*raw0 & 0x0000003fL) << 16;
         *cook   |= (*raw1 & 0x0003f000L) >> 4;
         *cook++ |= (*raw1 & 0x0000003fL);
     }
     usekey(dough);
     return;
 }

 void cpkey(into)
 register unsigned long *into;
 {
     register unsigned long *from, *endp;

     from = KnL, endp = &KnL[32];
     while( from < endp ) *into++ = *from++;
         return;
 }


 void des(inblock, outblock)
 unsigned char *inblock, *outblock;
 {
     unsigned long work[2];

     scrunch(inblock, work);
     desfunc(work, KnL);
     unscrun(work, outblock);
     return;
 }

 static void scrunch(outof, into)
 register unsigned char *outof;
 register unsigned long *into;
 {
     *into    = (*outof++ & 0xffL) << 24;
     *into   |= (*outof++ & 0xffL) << 16;
     *into   |= (*outof++ & 0xffL) << 8;
     *into++ |= (*outof++ & 0xffL);
     *into    = (*outof++ & 0xffL) << 24;
     *into   |= (*outof++ & 0xffL) << 16;
     *into   |= (*outof++ & 0xffL) << 8;
     *into   |= (*outof   & 0xffL);
     return;
 }

 static void unscrun(outof, into)
 register unsigned long *outof;
 register unsigned char *into;
 {
     *into++ = (unsigned char) ((*outof >> 24) & 0xffL);
     *into++ = (unsigned char) ((*outof >> 16) & 0xffL);
     *into++ = (unsigned char) ((*outof >>  8) & 0xffL);
     *into++ = (unsigned char)  (*outof++      & 0xffL);
     *into++ = (unsigned char) ((*outof >> 24) & 0xffL);
     *into++ = (unsigned char) ((*outof >> 16) & 0xffL);
     *into++ = (unsigned char) ((*outof >>  8) & 0xffL);
     *into   = (unsigned char)  (*outof        & 0xffL);
     return;
 }

 static unsigned long SP1[64] = {
     0x01010400L, 0x00000000L, 0x00010000L, 0x01010404L,
     0x01010004L, 0x00010404L, 0x00000004L, 0x00010000L,
     0x00000400L, 0x01010400L, 0x01010404L, 0x00000400L,
     0x01000404L, 0x01010004L, 0x01000000L, 0x00000004L,
     0x00000404L, 0x01000400L, 0x01000400L, 0x00010400L,
     0x00010400L, 0x01010000L, 0x01010000L, 0x01000404L,
     0x00010004L, 0x01000004L, 0x01000004L, 0x00010004L,
     0x00000000L, 0x00000404L, 0x00010404L, 0x01000000L,
     0x00010000L, 0x01010404L, 0x00000004L, 0x01010000L,
     0x01010400L, 0x01000000L, 0x01000000L, 0x00000400L,
     0x01010004L, 0x00010000L, 0x00010400L, 0x01000004L,
     0x00000400L, 0x00000004L, 0x01000404L, 0x00010404L,
     0x01010404L, 0x00010004L, 0x01010000L, 0x01000404L,
     0x01000004L, 0x00000404L, 0x00010404L, 0x01010400L,
     0x00000404L, 0x01000400L, 0x01000400L, 0x00000000L,
     0x00010004L, 0x00010400L, 0x00000000L, 0x01010004L };

 static unsigned long SP2[64] = {
     0x80108020L, 0x80008000L, 0x00008000L, 0x00108020L,
     0x00100000L, 0x00000020L, 0x80100020L, 0x80008020L,
     0x80000020L, 0x80108020L, 0x80108000L, 0x80000000L,
     0x80008000L, 0x00100000L, 0x00000020L, 0x80100020L,
     0x00108000L, 0x00100020L, 0x80008020L, 0x00000000L,
     0x80000000L, 0x00008000L, 0x00108020L, 0x80100000L,
     0x00100020L, 0x80000020L, 0x00000000L, 0x00108000L,
     0x00008020L, 0x80108000L, 0x80100000L, 0x00008020L,
     0x00000000L, 0x00108020L, 0x80100020L, 0x00100000L,
     0x80008020L, 0x80100000L, 0x80108000L, 0x00008000L,
     0x80100000L, 0x80008000L, 0x00000020L, 0x80108020L,
     0x00108020L, 0x00000020L, 0x00008000L, 0x80000000L,
     0x00008020L, 0x80108000L, 0x00100000L, 0x80000020L,
     0x00100020L, 0x80008020L, 0x80000020L, 0x00100020L,
     0x00108000L, 0x00000000L, 0x80008000L, 0x00008020L,
     0x80000000L, 0x80100020L, 0x80108020L, 0x00108000L };

 static unsigned long SP3[64] = {
     0x00000208L, 0x08020200L, 0x00000000L, 0x08020008L,
     0x08000200L, 0x00000000L, 0x00020208L, 0x08000200L,
     0x00020008L, 0x08000008L, 0x08000008L, 0x00020000L,
     0x08020208L, 0x00020008L, 0x08020000L, 0x00000208L,
     0x08000000L, 0x00000008L, 0x08020200L, 0x00000200L,
     0x00020200L, 0x08020000L, 0x08020008L, 0x00020208L,
     0x08000208L, 0x00020200L, 0x00020000L, 0x08000208L,
     0x00000008L, 0x08020208L, 0x00000200L, 0x08000000L,
     0x08020200L, 0x08000000L, 0x00020008L, 0x00000208L,
     0x00020000L, 0x08020200L, 0x08000200L, 0x00000000L,
     0x00000200L, 0x00020008L, 0x08020208L, 0x08000200L,
     0x08000008L, 0x00000200L, 0x00000000L, 0x08020008L,
     0x08000208L, 0x00020000L, 0x08000000L, 0x08020208L,
     0x00000008L, 0x00020208L, 0x00020200L, 0x08000008L,
     0x08020000L, 0x08000208L, 0x00000208L, 0x08020000L,
     0x00020208L, 0x00000008L, 0x08020008L, 0x00020200L };

 static unsigned long SP4[64] = {
     0x00802001L, 0x00002081L, 0x00002081L, 0x00000080L,
     0x00802080L, 0x00800081L, 0x00800001L, 0x00002001L,
     0x00000000L, 0x00802000L, 0x00802000L, 0x00802081L,
     0x00000081L, 0x00000000L, 0x00800080L, 0x00800001L,
     0x00000001L, 0x00002000L, 0x00800000L, 0x00802001L,
     0x00000080L, 0x00800000L, 0x00002001L, 0x00002080L,
     0x00800081L, 0x00000001L, 0x00002080L, 0x00800080L,
     0x00002000L, 0x00802080L, 0x00802081L, 0x00000081L,
     0x00800080L, 0x00800001L, 0x00802000L, 0x00802081L,
     0x00000081L, 0x00000000L, 0x00000000L, 0x00802000L,
     0x00002080L, 0x00800080L, 0x00800081L, 0x00000001L,
     0x00802001L, 0x00002081L, 0x00002081L, 0x00000080L,
     0x00802081L, 0x00000081L, 0x00000001L, 0x00002000L,
     0x00800001L, 0x00002001L, 0x00802080L, 0x00800081L,
     0x00002001L, 0x00002080L, 0x00800000L, 0x00802001L,
     0x00000080L, 0x00800000L, 0x00002000L, 0x00802080L };

 static unsigned long SP5[64] = {
     0x00000100L, 0x02080100L, 0x02080000L, 0x42000100L,
     0x00080000L, 0x00000100L, 0x40000000L, 0x02080000L,
     0x40080100L, 0x00080000L, 0x02000100L, 0x40080100L,
     0x42000100L, 0x42080000L, 0x00080100L, 0x40000000L,
     0x02000000L, 0x40080000L, 0x40080000L, 0x00000000L,
     0x40000100L, 0x42080100L, 0x42080100L, 0x02000100L,
     0x42080000L, 0x40000100L, 0x00000000L, 0x42000000L,
     0x02080100L, 0x02000000L, 0x42000000L, 0x00080100L,
     0x00080000L, 0x42000100L, 0x00000100L, 0x02000000L,
     0x40000000L, 0x02080000L, 0x42000100L, 0x40080100L,
     0x02000100L, 0x40000000L, 0x42080000L, 0x02080100L,
     0x40080100L, 0x00000100L, 0x02000000L, 0x42080000L,
     0x42080100L, 0x00080100L, 0x42000000L, 0x42080100L,
     0x02080000L, 0x00000000L, 0x40080000L, 0x42000000L,
     0x00080100L, 0x02000100L, 0x40000100L, 0x00080000L,
     0x00000000L, 0x40080000L, 0x02080100L, 0x40000100L };

 static unsigned long SP6[64] = {
     0x20000010L, 0x20400000L, 0x00004000L, 0x20404010L,
     0x20400000L, 0x00000010L, 0x20404010L, 0x00400000L,
     0x20004000L, 0x00404010L, 0x00400000L, 0x20000010L,
     0x00400010L, 0x20004000L, 0x20000000L, 0x00004010L,
     0x00000000L, 0x00400010L, 0x20004010L, 0x00004000L,
     0x00404000L, 0x20004010L, 0x00000010L, 0x20400010L,
     0x20400010L, 0x00000000L, 0x00404010L, 0x20404000L,
     0x00004010L, 0x00404000L, 0x20404000L, 0x20000000L,
     0x20004000L, 0x00000010L, 0x20400010L, 0x00404000L,
     0x20404010L, 0x00400000L, 0x00004010L, 0x20000010L,
     0x00400000L, 0x20004000L, 0x20000000L, 0x00004010L,
     0x20000010L, 0x20404010L, 0x00404000L, 0x20400000L,
     0x00404010L, 0x20404000L, 0x00000000L, 0x20400010L,
     0x00000010L, 0x00004000L, 0x20400000L, 0x00404010L,
     0x00004000L, 0x00400010L, 0x20004010L, 0x00000000L,
     0x20404000L, 0x20000000L, 0x00400010L, 0x20004010L };

 static unsigned long SP7[64] = {
     0x00200000L, 0x04200002L, 0x04000802L, 0x00000000L,
     0x00000800L, 0x04000802L, 0x00200802L, 0x04200800L,
     0x04200802L, 0x00200000L, 0x00000000L, 0x04000002L,
     0x00000002L, 0x04000000L, 0x04200002L, 0x00000802L,
     0x04000800L, 0x00200802L, 0x00200002L, 0x04000800L,
     0x04000002L, 0x04200000L, 0x04200800L, 0x00200002L,
     0x04200000L, 0x00000800L, 0x00000802L, 0x04200802L,
     0x00200800L, 0x00000002L, 0x04000000L, 0x00200800L,
     0x04000000L, 0x00200800L, 0x00200000L, 0x04000802L,
     0x04000802L, 0x04200002L, 0x04200002L, 0x00000002L,
     0x00200002L, 0x04000000L, 0x04000800L, 0x00200000L,
     0x04200800L, 0x00000802L, 0x00200802L, 0x04200800L,
     0x00000802L, 0x04000002L, 0x04200802L, 0x04200000L,
     0x00200800L, 0x00000000L, 0x00000002L, 0x04200802L,
     0x00000000L, 0x00200802L, 0x04200000L, 0x00000800L,
     0x04000002L, 0x04000800L, 0x00000800L, 0x00200002L };

 static unsigned long SP8[64] = {
     0x10001040L, 0x00001000L, 0x00040000L, 0x10041040L,
     0x10000000L, 0x10001040L, 0x00000040L, 0x10000000L,
     0x00040040L, 0x10040000L, 0x10041040L, 0x00041000L,
     0x10041000L, 0x00041040L, 0x00001000L, 0x00000040L,
     0x10040000L, 0x10000040L, 0x10001000L, 0x00001040L,
     0x00041000L, 0x00040040L, 0x10040040L, 0x10041000L,
     0x00001040L, 0x00000000L, 0x00000000L, 0x10040040L,
     0x10000040L, 0x10001000L, 0x00041040L, 0x00040000L,
     0x00041040L, 0x00040000L, 0x10041000L, 0x00001000L,
     0x00000040L, 0x10040040L, 0x00001000L, 0x00041040L,
     0x10001000L, 0x00000040L, 0x10000040L, 0x10040000L,
     0x10040040L, 0x10000000L, 0x00040000L, 0x10001040L,
     0x00000000L, 0x10041040L, 0x00040040L, 0x10000040L,
     0x10040000L, 0x10001000L, 0x10001040L, 0x00000000L,
     0x10041040L, 0x00041000L, 0x00041000L, 0x00001040L,
     0x00001040L, 0x00040040L, 0x10000000L, 0x10041000L };

 static void desfunc(block, keys)
 register unsigned long *block, *keys;
 {
     register unsigned long fval, work, right, leftt;
     register int round;

     leftt = block[0];
     right = block[1];
     work = ((leftt >> 4) ^ right) & 0x0f0f0f0fL;
     right ^= work;
     leftt ^= (work << 4);
     work = ((leftt >> 16) ^ right) & 0x0000ffffL;
     right ^= work;
     leftt ^= (work << 16);
     work = ((right >> 2) ^ leftt) & 0x33333333L;
     leftt ^= work;
     right ^= (work << 2);
     work = ((right >> 8) ^ leftt) & 0x00ff00ffL;
     leftt ^= work;
     right ^= (work << 8);
     right = ((right << 1) | ((right >> 31) & 1L)) & 0xffffffffL;
     work = (leftt ^ right) & 0xaaaaaaaaL;
     leftt ^= work;
     right ^= work;
     leftt = ((leftt << 1) | ((leftt >> 31) & 1L)) & 0xffffffffL;

     for( round = 0; round < 8; round++ ) {
         work  = (right << 28) | (right >> 4);
         work ^= *keys++;
         fval  = SP7[ work                & 0x3fL];
         fval |= SP5[(work >>  8) & 0x3fL];
         fval |= SP3[(work >> 16) & 0x3fL];
         fval |= SP1[(work >> 24) & 0x3fL];
         work  = right ^ *keys++;
         fval |= SP8[ work                & 0x3fL];
         fval |= SP6[(work >>  8) & 0x3fL];
         fval |= SP4[(work >> 16) & 0x3fL];
         fval |= SP2[(work >> 24) & 0x3fL];
         leftt ^= fval;
         work  = (leftt << 28) | (leftt >> 4);
         work ^= *keys++;
         fval  = SP7[ work                & 0x3fL];
         fval |= SP5[(work >>  8) & 0x3fL];
         fval |= SP3[(work >> 16) & 0x3fL];
         fval |= SP1[(work >> 24) & 0x3fL];
         work  = leftt ^ *keys++;
         fval |= SP8[ work                & 0x3fL];
         fval |= SP6[(work >>  8) & 0x3fL];
         fval |= SP4[(work >> 16) & 0x3fL];
         fval |= SP2[(work >> 24) & 0x3fL];
         right ^= fval;
     }

     right = (right << 31) | (right >> 1);
     work = (leftt ^ right) & 0xaaaaaaaaL;
     leftt ^= work;
     right ^= work;
     leftt = (leftt << 31) | (leftt >> 1);
     work = ((leftt >> 8) ^ right) & 0x00ff00ffL;
     right ^= work;
     leftt ^= (work << 8);
     work = ((leftt >> 2) ^ right) & 0x33333333L;
     right ^= work;
     leftt ^= (work << 2);
     work = ((right >> 16) ^ leftt) & 0x0000ffffL;
     leftt ^= work;
     right ^= (work << 16);
     work = ((right >> 4) ^ leftt) & 0x0f0f0f0fL;
     leftt ^= work;
     right ^= (work << 4);
     *block++ = right;
     *block = leftt;
     return;
 }

 /* Validation sets:
  *
  * Single-length key, single-length plaintext -
  * Key    : 0123 4567 89ab cdef
  * Plain  : 0123 4567 89ab cde7
  * Cipher : c957 4425 6a5e d31d
  *
  * Double-length key, single-length plaintext -
  * Key    : 0123 4567 89ab cdef fedc ba98 7654 3210
  * Plain  : 0123 4567 89ab cde7
  * Cipher : 7f1d 0a77 826b 8aff
  *
  * Double-length key, double-length plaintext -
  * Key    : 0123 4567 89ab cdef fedc ba98 7654 3210
  * Plain  : 0123 4567 89ab cdef 0123 4567 89ab cdff
  * Cipher : 27a0 8440 406a df60 278f 47cf 42d6 15d7
  *
  * Triple-length key, single-length plaintext -
  * Key    : 0123 4567 89ab cdef fedc ba98 7654 3210 89ab cdef 0123 4567
  * Plain  : 0123 4567 89ab cde7
  * Cipher : de0b 7c06 ae5e 0ed5
  *
  * Triple-length key, double-length plaintext -
  * Key    : 0123 4567 89ab cdef fedc ba98 7654 3210 89ab cdef 0123 4567
  * Plain  : 0123 4567 89ab cdef 0123 4567 89ab cdff
  * Cipher : ad0d 1b30 ac17 cf07 0ed1 1c63 81e4 4de5
  *
  * d3des V5.0a rwo 9208.07 18:44 Graven Imagery
  **********************************************************************/

 //  Copyright (C) 1999 AT&T Laboratories Cambridge. All Rights Reserved.
 //
 //  This file is part of the VNC system.
 //
 //  The VNC system is free software; you can redistribute it and/or modify
 //  it under the terms of the GNU General Public License as published by
 //  the Free Software Foundation; either version 2 of the License, or
 //  (at your option) any later version.
 //
 //  This program is distributed in the hope that it will be useful,
 //  but WITHOUT ANY WARRANTY; without even the implied warranty of
 //  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 //  GNU General Public License for more details.
 //
 //  You should have received a copy of the GNU General Public License
 //  along with this program; if not, write to the Free Software
 //  Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307,
 //  USA.
 //
 // TightVNC distribution homepage on the Web: http://www.tightvnc.com/
 //
 // If the source code for the VNC system is not available from the place whence
 // you received this file, check http://www.uk.research.att.com/vnc or contact
 // the authors on vnc@uk.research.att.com for information on obtaining it.


 /*
  * vncauth.c - Functions for VNC password management and authentication.
  */
 static const char *ID = "$Id: vncauth.c,v 1.5 2004/09/09 00:22:33 grmcdorman Exp $";

 #ifndef WIN32
 #include <stdio.h>
 #include <stdlib.h>
 #include <string.h>
 #include <sys/types.h>
 #include <sys/stat.h>
 #endif

 //  Copyright (C) 1999 AT&T Laboratories Cambridge. All Rights Reserved.
 //
 //  This file is part of the VNC system.
 //
 //  The VNC system is free software; you can redistribute it and/or modify
 //  it under the terms of the GNU General Public License as published by
 //  the Free Software Foundation; either version 2 of the License, or
 //  (at your option) any later version.
 //
 //  This program is distributed in the hope that it will be useful,
 //  but WITHOUT ANY WARRANTY; without even the implied warranty of
 //  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 //  GNU General Public License for more details.
 //
 //  You should have received a copy of the GNU General Public License
 //  along with this program; if not, write to the Free Software
 //  Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307,
 //  USA.
 //
 // TightVNC distribution homepage on the Web: http://www.tightvnc.com/
 //
 // If the source code for the VNC system is not available from the place
 // whence you received this file, check http://www.uk.research.att.com/vnc or contact
 // the authors on vnc@uk.research.att.com for information on obtaining it.


 /*
  * vncauth.h - describes the functions provided by the vncauth library.
  */

 #define MAXPWLEN 8
 #define CHALLENGESIZE 16

 #include <memory.h>

 /*
  * We use a fixed key to store passwords, since we assume that our local
  * file system is secure but nonetheless don't want to store passwords
  * as plaintext.
  */

 unsigned char fixedkey[8] = {23,82,107,6,35,78,88,7};


 /*
  * Encrypt a password and store it in a file.  Returns 0 if successful,
  * 1 if the file could not be written.
  */

 int
 vncEncryptAndStorePasswd(char *passwd, char *fname)
 {
     FILE *fp;
     int i;
     unsigned char encryptedPasswd[8];

     if ((fp = fopen(fname,"wb")) == NULL) return 1;

 #ifndef WIN32
     chmod(fname, S_IRUSR|S_IWUSR);
 #endif

     /* pad password with nulls */

     for (i = 0; i < 8; i++) {
         if (i < strlen(passwd)) {
             encryptedPasswd[i] = passwd[i];
         } else {
             encryptedPasswd[i] = 0;
         }
     }

     /* Do encryption in-place - this way we overwrite our copy of the plaintext
      password */

     deskey(fixedkey, EN0);
     des(encryptedPasswd, encryptedPasswd);

     for (i = 0; i < 8; i++) {
         putc(encryptedPasswd[i], fp);
     }

     fclose(fp);
     return 0;
 }

 /*
  * Decrypt a password from a file.  Returns a pointer to a newly allocated
  * string containing the password or a null pointer if the password could
  * not be retrieved for some reason.
  */

 char *
 vncDecryptPasswdFromFile(char *fname)
 {
     FILE *fp;
     int i, ch;
     unsigned char *passwd = (unsigned char *)malloc(9);

     if ((fp = fopen(fname,"r")) == NULL) return NULL;

     for (i = 0; i < 8; i++) {
         ch = getc(fp);
         if (ch == EOF) {
             fclose(fp);
             return NULL;
         }
         passwd[i] = ch;
     }

     fclose(fp);

     deskey(fixedkey, DE1);
     des(passwd, passwd);

     passwd[8] = 0;

     return (char *)passwd;
 }


 /*
  * Encrypt CHALLENGESIZE bytes in memory using a password.
  */

 void
 vncEncryptBytes(unsigned char *bytes, char *passwd)
 {
     unsigned char key[8];
     unsigned int i;

     /* key is simply password padded with nulls */

     for (i = 0; i < 8; i++) {
         if (i < strlen(passwd)) {
             key[i] = passwd[i];
         } else {
             key[i] = 0;
         }
     }

     deskey(key, EN0);

     for (i = 0; i < CHALLENGESIZE; i += 8) {
         des(bytes+i, bytes+i);
     }
 }


 /*
  * Encrypt a password into the specified space.
  * encryptedPasswd will be 8 bytes long - sufficient space
  *   should be allocated.
  */

 void
 vncEncryptPasswd( unsigned char *encryptedPasswd, char *passwd )
 {
     unsigned int i;

     /* pad password with nulls */
     for (i = 0; i < MAXPWLEN; i++) {
         if (i < strlen(passwd)) {
             encryptedPasswd[i] = passwd[i];
         } else {
             encryptedPasswd[i] = 0;
         }
     }

     /* Do encryption in-place - this way we overwrite our copy of the plaintext
      password */
     deskey(fixedkey, EN0);
     des(encryptedPasswd, encryptedPasswd);
 }


 /*
  * Decrypt a password.  Returns a pointer to a newly allocated
  * string containing the password or a null pointer if the password could
  * not be retrieved for some reason.
  */

 char *
 vncDecryptPasswd(const unsigned char *encryptedPasswd)
 {
     unsigned int i;
     unsigned char *passwd = (unsigned char *)malloc(MAXPWLEN+1);

     memcpy(passwd, encryptedPasswd, MAXPWLEN);

     for (i = 0; i < MAXPWLEN; i++) {
         passwd[i] = encryptedPasswd[i];
     }

     deskey(fixedkey, DE1);
     des(passwd, passwd);

     passwd[MAXPWLEN] = 0;

     return (char *)passwd;
 }
