import { createBrowserRouter } from 'react-router';
import { Layout } from './layout';
import { Home } from './screens/home';
import { Catalog } from './screens/catalog';
import { ProductDetail } from './screens/product-detail';
import { Favorites } from './screens/favorites';
import { FAQTopics } from './screens/faq-topics';
import { FAQArticles } from './screens/faq-articles';
import { FAQArticleDetail } from './screens/faq-article-detail';
import { Profile } from './screens/profile';
import { TermsOfService } from './screens/terms-of-service';
import { PrivacyPolicy } from './screens/privacy-policy';
import { NotFound } from './screens/not-found';

export const router = createBrowserRouter([
  {
    path: '/',
    Component: Layout,
    children: [
      { index: true, Component: Home },
      { path: 'catalog', Component: Catalog },
      { path: 'product/:id', Component: ProductDetail },
      { path: 'favorites', Component: Favorites },
      { path: 'faq', Component: FAQTopics },
      { path: 'faq/:topicId', Component: FAQArticles },
      { path: 'faq/:topicId/:articleId', Component: FAQArticleDetail },
      { path: 'profile', Component: Profile },
      { path: 'legal/terms', Component: TermsOfService },
      { path: 'legal/privacy', Component: PrivacyPolicy },
      { path: '*', Component: NotFound },
    ],
  },
]);
