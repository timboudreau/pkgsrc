# $NetBSD: Makefile,v 1.19 2019/04/26 14:12:49 maya Exp $

DISTNAME=	urlgrabber-3.1.0
PKGREVISION=	4
CATEGORIES=	www devel
MASTER_SITES=	http://urlgrabber.baseurl.org/download/

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://urlgrabber.baseurl.org/
COMMENT=	Python package that drastically simplifies the fetching of files

USE_LANGUAGES=	# none
PYTHON_VERSIONS_INCOMPATIBLE=	36 37 # not yet ported as of 3.1.0

.include "../../lang/python/application.mk"
.include "../../lang/python/distutils.mk"
.include "../../mk/bsd.pkg.mk"
