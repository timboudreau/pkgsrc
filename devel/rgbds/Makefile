# $NetBSD: Makefile,v 1.5 2019/03/04 13:18:05 nia Exp $

DISTNAME=	rgbds-0.3.8
CATEGORIES=	devel
MASTER_SITES=	${MASTER_SITE_GITHUB:=rednex/}
GITHUB_RELEASE=	v${PKGVERSION_NOREV}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://rednex.github.io/rgbds/
COMMENT=	Assembler/linker package for the Game Boy and Game Boy Color
LICENSE=	mit

EXTRACT_USING=	bsdtar
USE_TOOLS+=	gmake pkg-config lex bison

MAKE_FLAGS+=	mandir=${PREFIX}/${PKGMANDIR}

.include "../../graphics/png/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
