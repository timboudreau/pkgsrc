# $NetBSD: Makefile,v 1.8 2019/01/26 22:41:37 markd Exp $

DISTNAME=	hyphen-finnish
PKGNAME=	tex-${DISTNAME}-2018
TEXLIVE_REV=	49078

MAINTAINER=	minskim@NetBSD.org
COMMENT=	Hyphenation patterns for Finnish language
LICENSE=	lppl-1.3c

DEPENDS+=	tex-hyphen-base-[0-9]*:../../print/tex-hyphen-base
DEPENDS+=	tex-hyph-utf8-[0-9]*:../../print/tex-hyph-utf8

TEX_HYPHEN_DAT+=	finnish loadhyph-fi.tex
TEX_HYPHEN_DEF+=	\\addlanguage{finnish}{loadhyph-fi.tex}{}{2}{2}
TEX_HYPHEN_LUA+=	['finnish']={ \
			  loader='loadhyph-fi.tex', \
			  lefthyphenmin=2, \
			  righthyphenmin=2, \
			  synonyms={}, \
			  patterns='hyph-fi.pat.txt', \
			  hyphenation='', \
			},

.include "../../print/texlive/package.mk"
.include "../../mk/bsd.pkg.mk"
