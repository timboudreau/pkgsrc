# $NetBSD: Makefile,v 1.6 2019/04/16 18:41:20 bsiegert Exp $

DISTNAME=	minify-2.3.8
PKGNAME=	go-${DISTNAME}
PKGREVISION=	5
MASTER_SITES=	${MASTER_SITE_GITHUB:=tdewolff/}
CATEGORIES=	www
GITHUB_PROJECT=	minify
GITHUB_TAG=	v${PKGVERSION_NOREV}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://github.com/tdewolff/minify
COMMENT=	Go minifiers for web formats
LICENSE=	mit

GO_DIST_BASE=	${DISTNAME}
GO_SRCPATH=	github.com/tdewolff/minify

.include "../../www/go-parse/buildlink3.mk"
.include "../../devel/go-humanize/buildlink3.mk"
.include "../../devel/go-fsnotify/buildlink3.mk"
.include "../../devel/go-try/buildlink3.mk"
.include "../../devel/go-pflag/buildlink3.mk"
.include "../../devel/go-sys/buildlink3.mk"
.include "../../lang/go/go-package.mk"
.include "../../mk/bsd.pkg.mk"
