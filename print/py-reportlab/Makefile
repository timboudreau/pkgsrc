# $NetBSD: Makefile,v 1.70 2019/04/03 15:12:11 adam Exp $

DISTNAME=	reportlab-3.5.17
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	print python
MASTER_SITES=	${MASTER_SITE_PYPI:=r/reportlab/}

MAINTAINER=	joerg@NetBSD.org
HOMEPAGE=	http://www.reportlab.com/
COMMENT=	Powerful PDF-generating toolkit for Python
LICENSE=	modified-bsd

DEPENDS+=	${PYPKGPREFIX}-Pillow>=4.4.0:../../graphics/py-Pillow
DEPENDS+=	${PYPKGPREFIX}-expat-[0-9]*:../../textproc/py-expat
DEPENDS+=	${PYPKGPREFIX}-pyphen-[0-9]*:../../textproc/py-pyphen

.include "../../graphics/freetype2/buildlink3.mk"
.include "../../graphics/libart/buildlink3.mk"
.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
