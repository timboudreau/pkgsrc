# $NetBSD: Makefile,v 1.27 2018/08/10 09:00:36 adam Exp $

DISTNAME=	pandas-0.23.4
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	math graphics python
MASTER_SITES=	${MASTER_SITE_PYPI:=p/pandas/}

MAINTAINER=	bad@NetBSD.org
HOMEPAGE=	http://pandas.pydata.org/
COMMENT=	Python Data Analysis Library
LICENSE=	modified-bsd

DEPENDS+=	${PYPKGPREFIX}-bottleneck-[0-9]*:../../math/py-bottleneck
DEPENDS+=	${PYPKGPREFIX}-dateutil>=2.5.0:../../time/py-dateutil
DEPENDS+=	${PYPKGPREFIX}-matplotlib-[0-9]*:../../graphics/py-matplotlib
DEPENDS+=	${PYPKGPREFIX}-numexpr-[0-9]*:../../math/py-numexpr
DEPENDS+=	${PYPKGPREFIX}-pytz>=1.5:../../time/py-pytz
DEPENDS+=	${PYPKGPREFIX}-scipy>=0.7:../../math/py-scipy
DEPENDS+=	${PYPKGPREFIX}-sqlite3-[0-9]*:../../databases/py-sqlite3
DEPENDS+=	${PYPKGPREFIX}-tables>=2.2:../../math/py-tables
TEST_DEPENDS+=	${PYPKGPREFIX}-nose-[0-9]*:../../devel/py-nose
# 20 test failures as of 0.18.1, see
# https://github.com/pydata/pandas/issues/12337
# https://github.com/pydata/pandas/issues/14043

USE_LANGUAGES=	c c++

.include "../../lang/python/egg.mk"
BUILDLINK_API_DEPENDS.pynumpy+= ${PYPKGPREFIX}-numpy>=1.7.0
.include "../../math/py-numpy/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
